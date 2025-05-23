package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/olegdayo/omniconv"
)

func callMicroservice(method, endpoint string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("http://%s/%s", "gateway:8080", strings.TrimPrefix(endpoint, "/"))

	var body io.Reader = nil
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (s *Server) HasRuleUserProject(project_id int, user_id int64, rule string) bool {
	rawResp, err := callMicroservice("GET", fmt.Sprintf("auth/userProjectRole?user_id=%d&project_id=%d&role=%s", user_id, project_id, rule), nil)
	if err != nil {
		return false
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(rawResp, &parsed); err != nil {
		return false
	}
	hasRaw, ok := parsed["isRoleEnough"]
	has, ok2 := hasRaw.(bool)
	if !ok || !ok2 {
		return false
	}
	return has
}

func (s *Server) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims["user_id"]
	newproject, err := s.providerProject.CreateProject(context.Background(), ServerProject2ProviderProject(project))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}

	_, err = callMicroservice("POST", "auth/createProject", map[string]interface{}{"user_id": userID, "project_id": newproject.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}

	body, err := json.Marshal(ProviderProject2ServerProject(newproject))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getProjectsHanlder(w http.ResponseWriter, r *http.Request) {
	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))

	projects, err := s.providerProject.GetProjects(context.Background())
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawResp, err := callMicroservice("GET", fmt.Sprintf("auth/usersProjects?user_id=%d", userID), nil)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var projects1 map[string]interface{}
	if err := json.Unmarshal(rawResp, &projects1); err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ids_user, ok := projects1["project_ids"].([]interface{})
	if !ok {
		http.Error(w, "project_ids has unexpected type", http.StatusInternalServerError)
		return
	}

	idsSet := make(map[int64]struct{}, len(ids_user))
	for _, id := range ids_user {
		switch v := id.(type) {
		case float64:
			idsSet[int64(v)] = struct{}{}
		case int:
			idsSet[int64(v)] = struct{}{}
		case int64:
			idsSet[v] = struct{}{}
		}
	}

	type ProjectWithRole struct {
		Project interface{} `json:"project"`
		Role    string      `json:"role"`
	}

	var results []ProjectWithRole

	for _, p := range projects {
		if _, found := idsSet[int64(p.ID)]; !found {
			continue
		}

		var role string
		switch {
		case s.HasRuleUserProject(int(p.ID), userID, "owner"):
			role = "owner"
		case s.HasRuleUserProject(int(p.ID), userID, "editor"):
			role = "editor"
		case s.HasRuleUserProject(int(p.ID), userID, "viewer"):
			role = "viewer"
		default:
			continue
		}

		projectJSON := ProviderProject2ServerProject(p)

		results = append(results, ProjectWithRole{
			Project: projectJSON,
			Role:    role,
		})
	}

	body, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	err = s.providerProject.DeleteProject(context.Background(), project_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateProjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var project Project
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}

	err = s.providerProject.UpdateProject(context.Background(), project_id, ServerProject2ProviderProject(project))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetProjectGraphsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))

	if !s.HasRuleUserProject(projectID, userID, "viewer") {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Has not rule"))
		return
	}

	graphs, err := s.providerGraph.GetProjectGraphs(context.Background(), projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}

	converted := omniconv.ConvertSlice(graphs, ProviderGraph2ServerGraph)

	type response struct {
		CanEdit bool        `json:"can_edit"`
		Graphs  interface{} `json:"graphs"`
	}

	res := response{
		CanEdit: s.HasRuleUserProject(projectID, userID, "editor"),
		Graphs:  converted,
	}

	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) updateGraphHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	var graph Graph
	err = json.NewDecoder(r.Body).Decode(&graph)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	err = s.providerGraph.UpdateGraph(context.Background(), graph_id, ServerGraph2ProviderGraph(graph))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteGraphHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	err = s.providerGraph.DeleteGraph(context.Background(), graph_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateGraphServicesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var services []Service
	if err := json.NewDecoder(r.Body).Decode(&services); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	err = s.providerService.UpdateGraphServices(context.Background(), graph_id, omniconv.ConvertSlice(services, ServerService2ProviderService))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) createGraphHandler(w http.ResponseWriter, r *http.Request) {
	var graph Graph
	err := json.NewDecoder(r.Body).Decode(&graph)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	if !s.HasRuleUserProject(graph.ProjectID, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	newgraph, err := s.providerGraph.CreateGraph(context.Background(), ServerGraph2ProviderGraph(graph))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(ProviderGraph2ServerGraph(newgraph))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) updateGraphRelationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var relations []Relation
	if err := json.NewDecoder(r.Body).Decode(&relations); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	err = s.providerRelation.UpdateGraphRelations(context.Background(), graph_id, omniconv.ConvertSlice(relations, ServerRelation2ProviderRelation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getGraphServicesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "viewer") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	services, err := s.providerService.GetGraphServices(context.Background(), graph_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(omniconv.ConvertSlice(services, ProviderService2ServerService))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) getGraphRelationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "viewer") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	relations, err := s.providerRelation.GetGraphRelations(context.Background(), graph_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(omniconv.ConvertSlice(relations, ProviderRelation2ServerRelation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) createGraphServicesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var services []Service
	if err := json.NewDecoder(r.Body).Decode(&services); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	ids, err := s.providerService.CreateServices(context.Background(), graph_id, omniconv.ConvertSlice(services, ServerService2ProviderService))
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	body, err := json.Marshal(ids)
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	w.Write(body)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) createGraphRelationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var relations []Relation
	if err := json.NewDecoder(r.Body).Decode(&relations); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	err = s.providerRelation.CreateRelations(context.Background(), graph_id, omniconv.ConvertSlice(relations, ServerRelation2ProviderRelation))
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	graph_id, err := s.providerService.GetServiceGraph(context.Background(), service_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "viewer") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	service, err := s.providerService.GetService(context.Background(), service_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(ProviderService2ServerService(service))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	graph_id, err := s.providerService.GetServiceGraph(context.Background(), service_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	err = s.providerService.DeleteService(context.Background(), service_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	graph_id, err := s.providerService.GetServiceGraph(context.Background(), service_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var service Service
	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	err = s.providerService.UpdateService(context.Background(), service_id, ServerService2ProviderService(service))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) createServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), service.GraphID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	newservice, err := s.providerService.CreateService(context.Background(), ServerService2ProviderService(service))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(ProviderService2ServerService(newservice))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (s *Server) getRelationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relation_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	graph_id, err := s.providerRelation.GetRelationGraph(context.Background(), relation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "viewer") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	relation, err := s.providerRelation.GetRelation(context.Background(), relation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(ProviderRelation2ServerRelation(relation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) deleteRelationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relation_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	graph_id, err := s.providerRelation.GetRelationGraph(context.Background(), relation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	err = s.providerRelation.DeleteRelation(context.Background(), relation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) updateRelationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relation_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	graph_id, err := s.providerRelation.GetRelationGraph(context.Background(), relation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), graph_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "editor") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	var relation Relation
	err = json.NewDecoder(r.Body).Decode(&relation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	err = s.providerRelation.UpdateRelation(context.Background(), relation_id, ServerRelation2ProviderRelation(relation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) createRelationHandler(w http.ResponseWriter, r *http.Request) {
	var relation Relation
	err := json.NewDecoder(r.Body).Decode(&relation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	claims, ok := FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := int64(claims["user_id"].(float64))
	project_id, err := s.providerGraph.GetGraphProject(context.Background(), relation.GraphID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	if !s.HasRuleUserProject(project_id, userID, "viewer") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Has not rule"))
		return
	}

	newrelation, err := s.providerRelation.CreateRelation(context.Background(), ServerRelation2ProviderRelation(relation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(ProviderRelation2ServerRelation(newrelation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
