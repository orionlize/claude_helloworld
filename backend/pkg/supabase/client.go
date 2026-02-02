package supabase

import (
	"fmt"
	"os"
	"github.com/supabase/supabase-go"
)

// Client holds the Supabase client instance
var Client *supabase.Client

// InitSupabase initializes the Supabase client
func InitSupabase() error {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_ANON_KEY")

	if url == "" || key == "" {
		return fmt.Errorf("SUPABASE_URL and SUPABASE_ANON_KEY must be set")
	}

	var err error
	Client, err = supabase.NewClient(url, key, &supabase.ClientOptions{})
	if err != nil {
		return fmt.Errorf("failed to initialize supabase client: %w", err)
	}

	return nil
}

// GetSupabaseClient returns the Supabase client instance
func GetSupabaseClient() *supabase.Client {
	return Client
}
