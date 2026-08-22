package chat

import (
	"encoding/json"
	"net/http"
)

type HTTPHandler struct {
	Service *Service
}

type replyRequest struct {
	ConversationID string `json:"conversation_id"`
	ClientMode     bool   `json:"client_mode"`
	Message        string `json:"message"`
}

type replyResponse struct {
	ConversationID string `json:"conversation_id"`
	Reply          string `json:"reply"`
}

// ReplyHTTP is intentionally stateless for now. Persistence will be added with
// PostgreSQL after the core request/response flow is stable.
func (h *HTTPHandler) ReplyHTTP(w http.ResponseWriter, r *http.Request) {
	var input replyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if input.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	conversation := &Conversation{
		ID:         input.ConversationID,
		ClientMode: input.ClientMode,
	}
	if conversation.ID == "" {
		conversation.ID = "ephemeral"
	}

	reply, err := h.Service.Reply(r.Context(), conversation, input.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(replyResponse{
		ConversationID: conversation.ID,
		Reply:          reply,
	})
}
