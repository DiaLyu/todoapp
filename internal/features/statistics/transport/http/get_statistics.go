package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DiaLyu/todoapp/internal/core/domain"
	core_logger "github.com/DiaLyu/todoapp/internal/core/logger"
	core_http_request "github.com/DiaLyu/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/DiaLyu/todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TaskCreated               int      `json:"tasks_created"                 example:"50"`
	TaskCompleted             int      `json:"tasks_completed"               example:"10"`
	TaskCompletedRate         *float64 `json:"tasks_completed_rate"          example:"20"`
	TaskAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m30s"`
}

// GetTask        godoc
// @Summary       Получение статистики
// @Description   Просмотр статистики с опциональной фильтрацией по user_id автора задачи или временному промежутку
// @Tags          statistics
// @Produce       json
// @Param         user_id    query   int     false  "Фильтрация задач по ID автора"
// @Param         from       query   string  false  "Начало временного промежутка статистики (включительно)"
// @Param         to         query   string  false  "Конец временного промежутка статистики (включительно)"
// @Success       200 {object} GetStatisticsResponse   "Успешное получение статистики"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TaskAverageCompletionTime != nil {
		duration := statistics.TaskAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TaskCreated:               statistics.TaskCreated,
		TaskCompleted:             statistics.TaskCompleted,
		TaskCompletedRate:         statistics.TaskCompletedRate,
		TaskAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
