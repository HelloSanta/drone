// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

// +build !oss

package card

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drone/drone/handler/api/errors"
	"github.com/drone/drone/mock"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestHandleFindCardById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repos := mock.NewMockRepositoryStore(controller)
	repos.EXPECT().FindName(gomock.Any(), "octocat", "hello-world").Return(dummyRepo, nil)

	card := mock.NewMockCardStore(controller)
	card.EXPECT().FindById(gomock.Any(), dummyCard.Id).Return(dummyCard, nil)

	c := new(chi.Context)
	c.URLParams.Add("owner", "octocat")
	c.URLParams.Add("name", "hello-world")
	c.URLParams.Add("id", "1")

	in := new(bytes.Buffer)
	json.NewEncoder(in).Encode(dummyCreateCard)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", in)
	r = r.WithContext(
		context.WithValue(context.Background(), chi.RouteCtxKey, c),
	)

	HandleFindById(card, repos).ServeHTTP(w, r)
	if got, want := w.Code, http.StatusOK; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}
}

func TestHandleFindCardById_CardNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repos := mock.NewMockRepositoryStore(controller)
	repos.EXPECT().FindName(gomock.Any(), "octocat", "hello-world").Return(dummyRepo, nil)

	card := mock.NewMockCardStore(controller)
	card.EXPECT().FindById(gomock.Any(), dummyCard.Id).Return(nil, errors.ErrNotFound)

	c := new(chi.Context)
	c.URLParams.Add("owner", "octocat")
	c.URLParams.Add("name", "hello-world")
	c.URLParams.Add("id", "1")

	in := new(bytes.Buffer)
	json.NewEncoder(in).Encode(dummyCreateCard)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", in)
	r = r.WithContext(
		context.WithValue(context.Background(), chi.RouteCtxKey, c),
	)

	HandleFindById(card, repos).ServeHTTP(w, r)
	got, want := new(errors.Error), errors.ErrNotFound
	json.NewDecoder(w.Body).Decode(got)
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf(diff)
	}
}
