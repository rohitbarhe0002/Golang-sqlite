package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"my-go-project/internal/storage"
	"my-go-project/internal/types"
	"my-go-project/internal/utils/responce"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("craeting a student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			responce.WriteJson(w, http.StatusBadRequest, responce.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			responce.WriteJson(w, http.StatusBadRequest, responce.GeneralError(err))
			return
		}
		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			responce.WriteJson(w, http.StatusBadRequest, responce.ValidationError(validateErrs))
			return
		}
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user careated successfully ", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			responce.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		responce.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetbyID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting student ", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			responce.WriteJson(w, http.StatusBadRequest, responce.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Info("error getting user", slog.String("id", id))
			responce.WriteJson(w, http.StatusInternalServerError, responce.GeneralError(err))
			return
		}
		responce.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("getting all students ")
		students, err := storage.GetStudents()
		if err != nil {
			responce.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		responce.WriteJson(w, http.StatusOK, students)
	}
}
