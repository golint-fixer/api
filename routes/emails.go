package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

// EmailsListResponse contains the result of the EmailsList request.
type EmailsListResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Emails  *[]*models.Email `json:"emails,omitempty"`
}

// EmailsList sends a list of the emails in the inbox.
func EmailsList(c web.C, w http.ResponseWriter, r *http.Request) {
	// Fetch the current session from the database
	session := c.Env["token"].(*models.Token)

	// Parse the query
	var (
		query     = r.URL.Query()
		sortRaw   = query.Get("sort")
		offsetRaw = query.Get("offset")
		limitRaw  = query.Get("limit")
		thread    = query.Get("thread")
		sort      []string
		offset    int
		limit     int
	)

	if offsetRaw != "" {
		o, err := strconv.Atoi(offsetRaw)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error":  err,
				"offset": offset,
			}).Error("Invalid offset")

			utils.JSONResponse(w, 400, &EmailsListResponse{
				Success: false,
				Message: "Invalid offset",
			})
			return
		}
		offset = o
	}

	if limitRaw != "" {
		l, err := strconv.Atoi(limitRaw)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
				"limit": limit,
			}).Error("Invalid limit")

			utils.JSONResponse(w, 400, &EmailsListResponse{
				Success: false,
				Message: "Invalid limit",
			})
			return
		}
		limit = l
	}

	if sortRaw != "" {
		sort = strings.Split(sortRaw, ",")
	}

	// Get contacts from the database
	emails, err := env.Emails.List(session.Owner, sort, offset, limit, thread)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch emails")

		utils.JSONResponse(w, 500, &EmailsListResponse{
			Success: false,
			Message: "Internal error (code EM/LI/01)",
		})
		return
	}

	if offsetRaw != "" || limitRaw != "" {
		count, err := env.Emails.CountOwnedBy(session.Owner)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to count emails")

			utils.JSONResponse(w, 500, &EmailsListResponse{
				Success: false,
				Message: "Internal error (code EM/LI/02)",
			})
			return
		}
		w.Header().Set("X-Total-Count", strconv.Itoa(count))
	}

	utils.JSONResponse(w, 200, &EmailsListResponse{
		Success: true,
		Emails:  &emails,
	})

	// GET parameters:
	//   sort - split by commas, prefixes: - is desc, + is asc
	//   offset, limit - for pagination
	// Pagination ADDS X-Total-Count to the response!
}

type EmailsCreateRequest struct {
	To                  []string `json:"to"`
	CC                  []string `json:"cc"`
	BCC                 []string `json:"bcc"`
	ReplyTo             string   `json:"reply_to"`
	Thread              string   `json:"thread"`
	Subject             string   `json:"subject"`
	Body                string   `json:"body"`
	BodyVersionMajor    int      `json:"body_version_major"`
	BodyVersionMinor    int      `json:"body_version_minor"`
	Preview             string   `json:"preview"`
	PreviewVersionMajor int      `json:"preview_version_major"`
	PreviewVersionMinor int      `json:"preview_version_minor"`
	Encoding            string   `json:"encoding"`
	Attachments         []string `json:"attachments"`
	PGPFingerprints     []string `json:"pgp_fingerprints"`
}

// EmailsCreateResponse contains the result of the EmailsCreate request.
type EmailsCreateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Created []string `json:"created,omitempty"`
}

// EmailsCreate sends a new email
func EmailsCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var input EmailsCreateRequest
	err := utils.ParseRequest(r, &input)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Unable to decode a request")

		utils.JSONResponse(w, 400, &EmailsCreateResponse{
			Success: false,
			Message: "Invalid input format",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Ensure that the input data isn't empty
	if len(input.To) == 0 || input.Subject == "" || input.Body == "" {
		utils.JSONResponse(w, 400, &EmailsCreateResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Fetch the user object from the database
	account, err := env.Accounts.GetTokenOwner(c.Env["token"].(*models.Token))
	if err != nil {
		// The session refers to a non-existing user
		env.Log.WithFields(logrus.Fields{
			"id":    session.ID,
			"error": err.Error(),
		}).Warn("Valid session referred to a removed account")

		utils.JSONResponse(w, 410, &EmailsCreateResponse{
			Success: false,
			Message: "Account disabled",
		})
		return
	}

	// Create an email resource
	emailResource := models.MakeResource(session.Owner, input.Subject)

	// Get the "Sent" label's ID
	var label *models.Label
	err = env.Labels.WhereAndFetchOne(map[string]interface{}{
		"name":    "Sent",
		"builtin": true,
		"owner":   account.ID,
	}, &label)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"id":    account.ID,
			"error": err.Error(),
		}).Warn("Account has no sent label")

		utils.JSONResponse(w, 410, &EmailsCreateResponse{
			Success: false,
			Message: "Misconfigured account",
		})
		return
	}

	// Check if Thread is set
	if input.Thread != "" {
		// todo: make it an actual exists check to reduce lan bandwidth
		_, err := env.Threads.GetThread(input.Thread)
		if err != nil {
			env.Log.WithFields(logrus.Fields{
				"id":    input.Thread,
				"error": err.Error(),
			}).Warn("Cannot retrieve a thread")

			utils.JSONResponse(w, 400, &EmailsCreateResponse{
				Success: false,
				Message: "Invalid thread",
			})
			return
		}
	} else {
		thread := &models.Thread{
			Resource: models.MakeResource(account.ID, input.Subject),
			Emails:   []string{emailResource.ID},
			Labels:   []string{label.ID},
			Members:  append(append(input.To, input.CC...), input.BCC...),
			IsRead:   true,
		}

		err := env.Threads.Insert(thread)
		if err != nil {
			utils.JSONResponse(w, 500, &EmailsCreateResponse{
				Success: false,
				Message: "Unable to create a new thread",
			})

			env.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unable to create a new thread")
			return
		}

		input.Thread = thread.ID
	}

	// Create a new email struct
	email := &models.Email{
		Kind:        "sent",
		From:        []string{account.Name + "@" + env.Config.EmailDomain},
		To:          input.To,
		CC:          input.CC,
		BCC:         input.BCC,
		Resource:    emailResource,
		Attachments: input.Attachments,
		Thread:      input.Thread,
		Body: models.Encrypted{
			Encoding:        "json",
			PGPFingerprints: input.PGPFingerprints,
			Data:            input.Body,
			Schema:          "email_body",
			VersionMajor:    input.BodyVersionMajor,
			VersionMinor:    input.BodyVersionMinor,
		},
		Preview: models.Encrypted{
			Encoding:        "json",
			PGPFingerprints: input.PGPFingerprints,
			Data:            input.Preview,
			Schema:          "email_preview",
			VersionMajor:    input.PreviewVersionMajor,
			VersionMinor:    input.PreviewVersionMinor,
		},
		Status: "queued",
	}

	// Insert the email into the database
	if err := env.Emails.Insert(email); err != nil {
		utils.JSONResponse(w, 500, &EmailsCreateResponse{
			Success: false,
			Message: "internal server error - EM/CR/01",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not insert an email into the database")
		return
	}

	// Add a send request to the queue
	err = env.NATS.Publish("send", email.ID)
	if err != nil {
		utils.JSONResponse(w, 500, &EmailsCreateResponse{
			Success: false,
			Message: "internal server error - EM/CR/03",
		})

		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not publish an email send request")
		return
	}

	utils.JSONResponse(w, 201, &EmailsCreateResponse{
		Success: true,
		Created: []string{email.ID},
	})
}

// EmailsGetResponse contains the result of the EmailsGet request.
type EmailsGetResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Email   *models.Email `json:"email,omitempty"`
}

// EmailsGet responds with a single email message
func EmailsGet(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the email from the database
	email, err := env.Emails.GetEmail(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &EmailsGetResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if email.Owner != session.Owner {
		utils.JSONResponse(w, 404, &EmailsGetResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Write the email to the response
	utils.JSONResponse(w, 200, &EmailsGetResponse{
		Success: true,
		Email:   email,
	})
}

// EmailsDeleteResponse contains the result of the EmailsDelete request.
type EmailsDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EmailsDelete remvoes an email from the system
func EmailsDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the email from the database
	email, err := env.Emails.GetEmail(c.URLParams["id"])
	if err != nil {
		utils.JSONResponse(w, 404, &EmailsDeleteResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Fetch the current session from the middleware
	session := c.Env["token"].(*models.Token)

	// Check for ownership
	if email.Owner != session.Owner {
		utils.JSONResponse(w, 404, &EmailsDeleteResponse{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	// Perform the deletion
	err = env.Emails.DeleteID(c.URLParams["id"])
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.URLParams["id"],
		}).Error("Unable to delete a email")

		utils.JSONResponse(w, 500, &EmailsDeleteResponse{
			Success: false,
			Message: "Internal error (code EM/DE/01)",
		})
		return
	}

	// Write the email to the response
	utils.JSONResponse(w, 200, &EmailsDeleteResponse{
		Success: true,
		Message: "Email successfully removed",
	})
}
