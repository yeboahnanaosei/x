package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	myussd "github.com/yeboahnanaosei/x/testussd"
)

func handleResponse(ctx context.Context, user myussd.User, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("apiKey", os.Getenv("AFRICA_TALKING_API_KEY"))
	switch user.LastStage {
	case "":
		if err := userService.SetStage(ctx, user, "start"); err != nil {
			fmt.Println("error setting user stage:", err) // TODO: Log error
			w.WriteHeader(http.StatusInternalServerError)
			body := []byte("END we suffered an internal error")
			w.Header().Set("Content-Length", fmt.Sprint(len(body)))
			w.Write(body)
			return
		}
		w.Write([]byte("CON Welcome to MOE COVID-19 Tracker\n\n1. Report Case\n2. Request PPE"))
	case "start":
		if user.Input == "1" {
			w.Write([]byte("CON Enter victim name:\n"))
			userService.SetStage(ctx, user, "reportcase:name")
			return
		}
		if user.Input == "2" {
			w.Write([]byte("END Coming soon..."))
			userService.ClearStage(ctx, user)
			return
		}

		w.Write([]byte("END Invalid input. Please try again later"))
		userService.ClearStage(ctx, user)
	case "reportcase:name":
		c := myussd.Case{
			VictimName: user.Input,
			User:       user,
		}

		w.Write([]byte("CON Select victim's gender\n\n1. Female\n2. Male"))
		caseService.SetCase(ctx, c)
		userService.SetStage(ctx, user, "reportcase:gender")
	case "reportcase:gender":
		var gender string
		if user.Input == "1" {
			gender = "Female"
		} else if user.Input == "2" {
			gender = "Male"
		}
		w.Write([]byte("CON Select case type:\n\n1. Suspected\n2. Under Investigation\n3. Confirmed"))
		c, _ := caseService.GetCase(ctx, user)
		c.Gender = gender
		caseService.SetCase(ctx, c)
		userService.SetStage(ctx, user, "reportcase:status")
	case "reportcase:status":
		var caseStatus string
		if user.Input == "1" {
			caseStatus = "Suspected"
		}
		if user.Input == "2" {
			caseStatus = "Under Investigation"
		}
		if user.Input == "3" {
			caseStatus = "Confirmed"
		}

		c, _ := caseService.GetCase(ctx, user)
		c.Status = caseStatus
		output := fmt.Sprintf("CON Please confirm\n%s\n1. Confirm\n2. Cancel", c)
		w.Write([]byte(output))
		caseService.SetCase(ctx, c)
		userService.SetStage(ctx, user, "reportcase:confirm")
	case "reportcase:confirm":
		c, _ := caseService.GetCase(ctx, user)
		if user.Input == "1" {
			w.Write([]byte("END Case reported successfully. Thank you"))
			userService.ClearStage(ctx, user)
			caseService.ClearCase(ctx, c)
		} else if user.Input == "2" {
			w.Write([]byte("END Case cancelled. Thank you"))
			userService.ClearStage(ctx, user)
			caseService.ClearCase(ctx, c)
		} else {
			w.Write([]byte("END Unknown option"))
		}
		userService.ClearStage(ctx, user)
		caseService.ClearCase(ctx, c)
	case "requestppe":
		w.Write([]byte("END Coming soon..."))
	default:
		fmt.Println("unknown user stage:", user.LastStage) // TODO: Log error
		w.WriteHeader(http.StatusInternalServerError)
		body := []byte("END we suffered an internal error")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.Write(body)
	}
}
