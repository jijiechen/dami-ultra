package behavior_tests

import (
	"fmt"

	"github.com/stretchr/testify/require"
	"github.com/yuchanns/kong-exercise-microservices/internal/business"
)

func testServiceCreateAndGet(assert *require.Assertions, b *reqBuilder) {
	name := generateRandomString(4)
	desc := generateRandomString(10)

	_ = b.Record("POST", "/service/", &business.ServiceVersion{
		Service: business.Service{
			Name:        name,
			Description: desc,
		},
		Version: business.Version{
			Host:     "test.com",
			Port:     443,
			Path:     "/api",
			Protocol: "https",
		},
	})

	buf := b.Record("GET", "/service/list", business.ListOpt{
		Keyword: name,
	})

	result := parse[business.ListRes](assert, buf)

	assert.Greater(result.Total, int64(0))
	var id int64
	for _, v := range result.Items {
		if v.Name == name {
			id = v.ID
			assert.Equal(desc, v.Description)
			break
		}
	}
	assert.True(id > 0)

	buf = b.Record("GET", "/service/"+fmt.Sprintf("%d", id), nil)
	detail := parse[business.ServiceDetail](assert, buf)
	assert.Equal(name, detail.Name)
	assert.Equal(1, len(detail.Versions))
}

func testServiceList(assert *require.Assertions, b *reqBuilder) {
	prefix := generateRandomString(2)
	total := 20
	for i := 0; i < total; i++ {
		name := prefix + generateRandomString(2)
		desc := generateRandomString(10)
		_ = b.Record("POST", "/service/", &business.ServiceVersion{
			Service: business.Service{
				Name:        name,
				Description: desc,
			},
			Version: business.Version{
				Host:     "test.com",
				Port:     443,
				Path:     "/api",
				Protocol: "https",
			},
		})
	}
	buf := b.Record("GET", "/service/list", business.ListOpt{
		Keyword: prefix,
		Pager:   &business.Pager{Page: 1, PageSize: int32(total)},
	})
	result := parse[business.ListRes](assert, buf)

	assert.Equal(int64(total), result.Total)
	assert.Equal(int64(total), int64(len(result.Items)))

	buf = b.Record("GET", "/service/list", business.ListOpt{
		Keyword: prefix,
		Pager:   &business.Pager{Page: 1, PageSize: 5},
	})
	result = parse[business.ListRes](assert, buf)
	assert.Equal(int64(total), result.Total)
	assert.Equal(int64(5), int64(len(result.Items)))
	for _, v := range result.Items {
		assert.Equal(int32(1), v.VersionCnt)
	}
}

func testServiceUpdate(assert *require.Assertions, b *reqBuilder) {
	name := generateRandomString(4)
	desc := generateRandomString(10)

	hosts := []string{"test.com"}

	_ = b.Record("POST", "/service/", &business.ServiceVersion{
		Service: business.Service{
			Name:        name,
			Description: desc,
		},
		Version: business.Version{
			Host:     "test.com",
			Port:     443,
			Path:     "/api",
			Protocol: "https",
		},
	})

	buf := b.Record("GET", "/service/list", business.ListOpt{
		Keyword: name,
	})

	result := parse[business.ListRes](assert, buf)

	assert.Greater(result.Total, int64(0))
	var id int64
	for _, v := range result.Items {
		if v.Name == name {
			id = v.ID
			assert.Equal(desc, v.Description)
			break
		}
	}
	assert.True(id > 0)

	total := generateRandomInt(1, 10)

	for i := int64(0); i < total; i++ {
		host := generateRandomString(4) + ".com"
		hosts = append(hosts, host)
		_ = b.Record("PUT", "/service/", &business.ServiceVersion{
			Service: business.Service{
				Name:        name,
				Description: desc,
				ID:          id,
			},
			Version: business.Version{
				Host:     host,
				Port:     443,
				Path:     "/api",
				Protocol: "http",
			},
		})
	}

	buf = b.Record("GET", "/service/"+fmt.Sprintf("%d", id), nil)
	detail := parse[business.ServiceDetail](assert, buf)
	assert.Equal(name, detail.Name)
	assert.Equal(total+1, int64(len(detail.Versions)))
	for _, v := range detail.Versions {
		assert.Contains(hosts, v.Host)
	}
}
