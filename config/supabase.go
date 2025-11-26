package config

type Supabase interface{
    NewSupabaseConfig(url string, key string, bucket string)
}

type SupabaseConfig struct {
    URL        string
    ServiceKey string
    BucketName string
}

func NewSupabaseConfig(url string, key string, bucket string) SupabaseConfig {
    return SupabaseConfig{
        URL:        url,
        ServiceKey: key,
        BucketName: bucket,
    }
}