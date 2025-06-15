package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/olegdayo/omniconv"
)

func (s *Server) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
	}

	newproject, err := s.providerProject.CreateProject(r.Context(), ServerProject2ProviderProject(project))
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
	projects, err := s.providerProject.GetProjects(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(omniconv.ConvertSlice(projects, ProviderProject2ServerProject))
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

	err = s.providerProject.DeleteProject(r.Context(), project_id)
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

	err = s.providerProject.UpdateProject(r.Context(), project_id, ServerProject2ProviderProject(project))
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
	graphs, err := s.providerGraph.GetProjectGraphs(r.Context(), project_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(omniconv.ConvertSlice(graphs, ProviderGraph2ServerGraph))
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

	err = s.providerGraph.UpdateGraph(r.Context(), graph_id, ServerGraph2ProviderGraph(graph))
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

	err = s.providerGraph.DeleteGraph(r.Context(), graph_id)
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
	var services []Service
	if err := json.NewDecoder(r.Body).Decode(&services); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	err = s.providerService.UpdateGraphServices(r.Context(), graph_id, omniconv.ConvertSlice(services, ServerService2ProviderService))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateGraphRelationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	var relations []Relation
	if err := json.NewDecoder(r.Body).Decode(&relations); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	err = s.providerRelation.UpdateGraphRelations(r.Context(), graph_id, omniconv.ConvertSlice(relations, ServerRelation2ProviderRelation))
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

	newgraph, err := s.providerGraph.CreateGraph(r.Context(), ServerGraph2ProviderGraph(graph))
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

func (s *Server) getGraphServicesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	graph_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	services, err := s.providerService.GetGraphServices(r.Context(), graph_id)
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

	relations, err := s.providerRelation.GetGraphRelations(r.Context(), graph_id)
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
	var services []Service
	if err := json.NewDecoder(r.Body).Decode(&services); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	ids, err := s.providerService.CreateServices(r.Context(), graph_id, omniconv.ConvertSlice(services, ServerService2ProviderService))
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
	var relations []Relation
	if err := json.NewDecoder(r.Body).Decode(&relations); err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	err = s.providerRelation.CreateRelations(r.Context(), graph_id, omniconv.ConvertSlice(relations, ServerRelation2ProviderRelation))
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

	service, err := s.providerService.GetService(r.Context(), service_id)
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

	err = s.providerService.DeleteService(r.Context(), service_id)
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

	err = s.providerService.UpdateService(r.Context(), service_id, ServerService2ProviderService(service))
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

	newservice, err := s.providerService.CreateService(r.Context(), ServerService2ProviderService(service))
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

	relation, err := s.providerRelation.GetRelation(r.Context(), relation_id)
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

	err = s.providerRelation.DeleteRelation(r.Context(), relation_id)
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

	err = s.providerRelation.UpdateRelation(r.Context(), relation_id, ServerRelation2ProviderRelation(relation))
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

	newrelation, err := s.providerRelation.CreateRelation(r.Context(), ServerRelation2ProviderRelation(relation))
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
