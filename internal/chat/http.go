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
	OwnerID        string `json:"owner_id"`
	ContactID      string `json:"contact_id"`
	ClientMode     bool   `json:"client_mode"`
	Message        string `json:"message"`
}

type replyResponse struct {
	ConversationID string `json:"conversation_id"`
	Decision       string `json:"decision"`
	Reply          string `json:"reply,omitempty"`
}

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

	conversation := NewConversation(input.ConversationID, input.OwnerID, input.ContactID, input.ClientMode)
	if h.Service.Repository != nil && input.ConversationID != "" {
		stored, err := h.Service.Repository.Get(r.Context(), input.ConversationID)
		if err == nil {
			conversation = stored
		} else if err != ErrConversationNotFound {
			http.Error(w, "failed to load conversation", http.StatusInternalServerError)
			return
		}
	}

	result, err := h.Service.Reply(r.Context(), &conversation, input.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(replyResponse{
		ConversationID: conversation.ID,
		Decision:       string(result.Decision),
		Reply:          result.Response,
	})
}
