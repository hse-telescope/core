package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	newproject, err := s.providerProject.CreateProject(context.Background(), ServerProject2ProviderProject(project))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(newproject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getProjectsHanlder(w http.ResponseWriter, r *http.Request) {
	projects, err := s.providerProject.GetProjects(context.Background())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(projects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
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

	var project Project
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
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
	project_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	graphs, err := s.providerGraph.GetProjectGraphs(context.Background(), project_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(graphs)
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

	err = s.providerGraph.DeleteGraph(context.Background(), graph_id)
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

	newgraph, err := s.providerGraph.CreateGraph(context.Background(), ServerGraph2ProviderGraph(graph))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(newgraph)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) getGraphServicesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	services, err := s.providerService.GetGraphServices(context.Background(), graph_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(services)
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

	relations, err := s.providerRelation.GetGraphRelations(context.Background(), graph_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(relations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	service, err := s.providerService.GetService(context.Background(), service_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(service)
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

	newservice, err := s.providerService.CreateService(context.Background(), ServerService2ProviderService(service))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(newservice)
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

	relation, err := s.providerRelation.GetRelation(context.Background(), relation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(relation)
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

	newrelation, err := s.providerRelation.CreateRelation(context.Background(), ServerRelation2ProviderRelation(relation))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(newrelation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
