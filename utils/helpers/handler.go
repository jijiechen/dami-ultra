package helpers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BuildHandler is a helper function that wraps a function that takes a context and a request
// and returns a response and an error.
// It returns a gin.HandlerFunc that calls the function and writes the response to the context.
// If the function returns an error, it writes the error to the context with a 500 status code.
// This is useful for seperating the abstract logic of a handler from the gin specifics.
func BuildHandler[T any, R any](fn func(ctx context.Context, r R) (T, error)) gin.HandlerFunc {
	v := validator.New()
	return func(ctx *gin.Context) {
		var r R
		if err := ctx.ShouldBind(&r); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := v.Struct(r)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := fn(ctx.Request.Context(), r)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func BuildHandlerUri[T any, R any](fn func(ctx context.Context, r R) (T, error)) gin.HandlerFunc {
	v := validator.New()
	return func(ctx *gin.Context) {
		var r R
		if err := ctx.ShouldBindUri(&r); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := v.Struct(r)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := fn(ctx.Request.Context(), r)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": data})
	}
}
