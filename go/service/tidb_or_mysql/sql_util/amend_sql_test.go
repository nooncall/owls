package sql_util

import "testing"

func TestAddLimit(t *testing.T) {
	testData := []struct {
		source, expect string
	}{
		{
			"select * from cluster;  ",
			"select * from cluster limit 10",
		},
		{
			"select * from cluster limit 2",
			"select * from cluster limit 2",
		},
		{
			"select * from cluster",
			"select * from cluster limit 10",
		},

		{
			"select * from owl_clusters where id in (select id from owl_clusters limit 2)",
			"select * from owl_clusters where id in (select id from owl_clusters limit 2) limit 10",
		},
		{
			"select * from owl_clusters where id in (select id from owl_clusters)",
			"select * from owl_clusters where id in (select id from owl_clusters) limit 10",
		},
		{
			"select * from owl_clusters where id in (select id from owl_clusters) limit 5",
			"select * from owl_clusters where id in (select id from owl_clusters) limit 5",
		},
	}

	for i, v := range testData {
		if v.expect != AddLimitToSelect(v.source) {
			t.Log("failed at :", testData[i])
			t.FailNow()
		}
	}
}
