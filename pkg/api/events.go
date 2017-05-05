/*******************************************************************************
*
* Copyright 2017 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package api

import (
	"net/http"

	"github.com/sapcc/hermes/pkg/hermes"
)

//ListEvents handles GET /v1/events.
func (p *v1Provider) ListEvents(w http.ResponseWriter, r *http.Request) {
	if !p.CheckToken(r).Require(w, "event:list") {
		return
	}

	events, err := hermes.GetEvents(nil, nil, nil)
	if ReturnError(w, err) {
		return
	}

	ReturnJSON(w, 200, map[string]interface{}{"events": events})
}

//GetEvent handles GET /v1/events/:event_id.
func (p *v1Provider) GetEvent(w http.ResponseWriter, r *http.Request) {
	if !p.CheckToken(r).Require(w, "event:show") {
		return
	}

	events, err := hermes.GetEvents(nil, nil, nil)
	if ReturnError(w, err) {
		return
	}

	ReturnJSON(w, 200, map[string]interface{}{"event": events[0]})
}