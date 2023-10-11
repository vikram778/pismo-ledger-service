package server

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"pismo-ledger-service/errs"
	"pismo-ledger-service/pkg/models"
	"strconv"
)

func (h *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account

	err := h.GetParams(&account, r)
	if err != nil {
		h.FormatException(w, err)
		return
	}

	acct, err := h.svc.GetAccountByDocument(context.Background(), account.DocumentNumber)
	if err != nil {
		h.FormatException(w, err)
		return
	}

	if acct.AccountID != 0 {
		h.FormatException(w, errs.ErrorAccountExist)
		return

	}

	err = h.svc.CreateAccount(context.Background(), &account)
	if err != nil {
		h.FormatException(w, err)
		return
	}

	h.JSON(w, http.StatusOK, "account created successfully!!")
	return
}

func (h *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["account_id"]

	acctId, _ := strconv.Atoi(id)

	acct, err := h.svc.GetAccount(context.Background(), int64(acctId))
	if err != nil {
		h.FormatException(w, errs.ErrorAccountNotExist)
		return
	}
	h.JSON(w, http.StatusOK, acct)
	return
}

func (h *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := h.GetParams(&transaction, r)
	if err != nil {
		h.FormatException(w, err)
		return
	}

	_, err = h.svc.GetAccount(context.Background(), transaction.AccountID)
	if err != nil {
		h.FormatException(w, errs.ErrorAccountNotExist)
		return
	}

	_, err = h.svc.GetOperationType(context.Background(), transaction.OperationTypeID)
	if err != nil {
		h.FormatException(w, errs.ErrorIncorrectOperationType)
		return
	}

	if transaction.OperationTypeID != 4 {
		transaction.Amount *= -1
	}

	err = h.svc.CreateTransaction(context.Background(), &transaction)
	if err != nil {
		h.FormatException(w, err)
		return
	}

	h.JSON(w, http.StatusOK, "success")
	return
}
