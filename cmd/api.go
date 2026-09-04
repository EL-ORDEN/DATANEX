package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"datanex/internal/api"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Call REST APIs",
}

var apiGetCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Make a GET request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return callAPI(cmd, "GET", args[0], nil)
	},
}

var apiPostCmd = &cobra.Command{
	Use:   "post [url]",
	Short: "Make a POST request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := cmd.Flags().GetString("data")
		var data any
		if payload != "" {
			if err := json.Unmarshal([]byte(payload), &data); err != nil {
				return fmt.Errorf("decode JSON body: %w", err)
			}
		}
		return callAPI(cmd, "POST", args[0], data)
	},
}

var apiPutCmd = &cobra.Command{
	Use:   "put [url]",
	Short: "Make a PUT request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := cmd.Flags().GetString("data")
		var data any
		if payload != "" {
			if err := json.Unmarshal([]byte(payload), &data); err != nil {
				return fmt.Errorf("decode JSON body: %w", err)
			}
		}
		return callAPI(cmd, "PUT", args[0], data)
	},
}

var apiPatchCmd = &cobra.Command{
	Use:   "patch [url]",
	Short: "Make a PATCH request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := cmd.Flags().GetString("data")
		var data any
		if payload != "" {
			if err := json.Unmarshal([]byte(payload), &data); err != nil {
				return fmt.Errorf("decode JSON body: %w", err)
			}
		}
		return callAPI(cmd, "PATCH", args[0], data)
	},
}

var apiDeleteCmd = &cobra.Command{
	Use:   "delete [url]",
	Short: "Make a DELETE request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return callAPI(cmd, "DELETE", args[0], nil)
	},
}

func callAPI(cmd *cobra.Command, method, url string, body any) error {
	headerMap := make(map[string]string)
	for _, h := range os.Environ() {
		if strings.HasPrefix(h, "DATANEX_HEADER_") {
			parts := strings.SplitN(h, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "DATANEX_HEADER_")
				headerMap[key] = parts[1]
			}
		}
	}
	if headerValue, _ := cmd.Flags().GetString("header"); headerValue != "" {
		for _, entry := range strings.Split(headerValue, ",") {
			parts := strings.SplitN(entry, ":", 2)
			if len(parts) == 2 {
				headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	bearerToken, _ := cmd.Flags().GetString("token")
	client := api.NewClient()
	resp, err := client.Do(api.Request{
		Method:      method,
		URL:         url,
		Headers:     headerMap,
		Body:        body,
		BearerToken: bearerToken,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n\n", method, url)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Time: %s\n\n", resp.Time.Round(time.Millisecond))
	if len(resp.Body) == 0 {
		fmt.Println("<empty response>")
		return nil
	}
	var pretty any
	if err := json.Unmarshal(resp.Body, &pretty); err == nil {
		indented, err := json.MarshalIndent(pretty, "", "  ")
		if err == nil {
			fmt.Println(string(indented))
			return nil
		}
	}
	fmt.Println(string(resp.Body))
	return nil
}

func init() {
	apiCmd.AddCommand(apiGetCmd, apiPostCmd, apiPutCmd, apiPatchCmd, apiDeleteCmd)
	rootCmd.AddCommand(apiCmd)
	apiCmd.PersistentFlags().String("header", "", "Header in key:value format, optionally comma-separated")
	apiCmd.PersistentFlags().String("token", "", "Bearer token for Authorization header")
	apiPostCmd.Flags().String("data", "", "JSON body for POST requests")
	apiPutCmd.Flags().String("data", "", "JSON body for PUT requests")
	apiPatchCmd.Flags().String("data", "", "JSON body for PATCH requests")
}
