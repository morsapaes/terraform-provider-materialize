package resources

import (
	"context"
	"testing"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/testhelpers"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestResourceClusterReplicaCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                          "replica",
		"cluster_name":                  "cluster",
		"size":                          "small",
		"availability_zone":             "use1-az1",
		"introspection_interval":        "10s",
		"introspection_debugging":       true,
		"idle_arrangement_merge_effort": 100,
	}
	d := schema.TestResourceDataRaw(t, ClusterReplica().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		// Create
		mock.ExpectExec(
			`CREATE CLUSTER REPLICA "cluster"."replica" SIZE = 'small', AVAILABILITY ZONE = 'use1-az1', INTROSPECTION INTERVAL = '10s', INTROSPECTION DEBUGGING = TRUE, IDLE ARRANGEMENT MERGE EFFORT = 100;`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Query Id
		ir := mock.NewRows([]string{"id", "replica_name", "cluster_name", "size", "availability_zone"}).
			AddRow("u1", "replica", "cluster", "small", "use1-az2")
		mock.ExpectQuery(`
			SELECT
				mz_cluster_replicas.id,
				mz_cluster_replicas.name AS replica_name,
				mz_clusters.name AS cluster_name,
				mz_cluster_replicas.size,
				mz_cluster_replicas.availability_zone
			FROM mz_cluster_replicas
			JOIN mz_clusters
				ON mz_cluster_replicas.cluster_id = mz_clusters.id
			WHERE mz_cluster_replicas.name = 'replica'
			AND mz_clusters.name = 'cluster';`).WillReturnRows(ir)

		// Query Params
		ip := mock.NewRows([]string{"id", "replica_name", "cluster_name", "size", "availability_zone"}).
			AddRow("u1", "replica", "cluster", "small", "use1-az2")
		mock.ExpectQuery(`
			SELECT
				mz_cluster_replicas.id,
				mz_cluster_replicas.name AS replica_name,
				mz_clusters.name AS cluster_name,
				mz_cluster_replicas.size,
				mz_cluster_replicas.availability_zone
			FROM mz_cluster_replicas
			JOIN mz_clusters
				ON mz_cluster_replicas.cluster_id = mz_clusters.id
			WHERE mz_cluster_replicas.id = 'u1';`).WillReturnRows(ip)

		if err := clusterReplicaCreate(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}

	})

}

func TestResourceClusterReplicaDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":         "replica",
		"cluster_name": "cluster",
	}
	d := schema.TestResourceDataRaw(t, ClusterReplica().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP CLUSTER REPLICA "cluster"."replica";`).WillReturnResult(sqlmock.NewResult(1, 1))

		if err := clusterReplicaDelete(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}
