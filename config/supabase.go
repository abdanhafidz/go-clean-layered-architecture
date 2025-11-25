package config

type SupabaseConfig struct {
    URL        string
    ServiceKey string
    BucketName string
}

func NewSupabaseConfig(url, key, bucket string) SupabaseConfig {
    return SupabaseConfig{
        URL:        url,
        ServiceKey: key,
        BucketName: bucket,
    }
}