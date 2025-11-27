package config

type SupabaseConfig interface {
    GetURL() string
    GetServiceKey() string
    GetBucketName() string
}

type supabaseConfig struct {
    url        string
    serviceKey string
    bucketName string
}

func NewSupabaseConfig(url string, key string, bucket string) SupabaseConfig {
    return &supabaseConfig{
        url:        url,
        serviceKey: key,
        bucketName: bucket,
    }
}

func (c *supabaseConfig) GetURL() string { return c.url }
func (c *supabaseConfig) GetServiceKey() string { return c.serviceKey }
func (c *supabaseConfig) GetBucketName() string { return c.bucketName }
