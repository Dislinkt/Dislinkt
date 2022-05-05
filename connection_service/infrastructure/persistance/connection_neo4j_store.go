package persistance

import (
	"fmt"
	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"time"
)

type ConnectionDBStore struct {
	connectionDB *neo4j.Driver
}

func NewConnectionDBStore(client *neo4j.Driver) domain.ConnectionStore {
	return &ConnectionDBStore{
		connectionDB: client,
	}
}

func (store *ConnectionDBStore) CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(baseUserUuid, tx) && checkIfUserExist(connectUserUuid, tx) {

			records, err := tx.Run("MATCH (n:UserNode { uid: $uid}) RETURN n.status", map[string]interface{}{
				"uid": baseUserUuid,
			})
			connection_status := ""
			if err != nil {
				return nil, err
			}
			record, err := records.Single()
			if err != nil {
				return nil, err
			}

			status := record.Values[0].(string)
			dateNow := time.Now().Local().Unix()
			if status == "PRIVATE" {
				connection_status = "REQUEST_SENT"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $connect_user_uuid  MATCH (u2:UserNode) WHERE u2.uid = $base_user_uuid CREATE (u1)-[r1:CONNECTION {status: $status, date: $date}]->(u2)", map[string]interface{}{
					"connect_user_uuid": connectUserUuid,
					"base_user_uuid":    baseUserUuid,
					"status":            "REQUEST_SENT",
					"date":              dateNow,
				})

				if err != nil {
					return nil, err
				}

			} else {
				connection_status = "CONNECTED"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = &connect_user_uuid  MATCH (u2:UserNode) WHERE u2.uid = &base_user_uuid CREATE (u1)-[r1:CONNECTION {status: &status, date: &date}]->(u2) CREATE (u2)-[r1:CONNECTION {status: &status, date: &date}]->(u1)", map[string]interface{}{
					"connect_user_uuid": connectUserUuid,
					"base_user_uuid":    baseUserUuid,
					"status":            "CONNECTED",
					"date":              dateNow,
				})

				if err != nil {
					return nil, err
				}
			}

			fmt.Println(status)

			// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
			return &pb.NewConnectionResponse{
				BaseUserUUID:       baseUserUuid,
				ConnectUserUUID:    connectUserUuid,
				ConnectionResponse: connection_status,
			}, nil
		} else {
			return &pb.NewConnectionResponse{
				BaseUserUUID:       baseUserUuid,
				ConnectUserUUID:    connectUserUuid,
				ConnectionResponse: "Connection refused",
			}, nil
		}

	})

	if err != nil {
		return nil, err
	}

	return result.(*pb.NewConnectionResponse), nil
}

func (store *ConnectionDBStore) Register(userNode *domain.UserNode) (*domain.UserNode, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(userNode.UserUID, tx) {
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		records, err := tx.Run("CREATE (n:UserNode { uid: $uid, status: $status  }) RETURN n.uid, n.status", map[string]interface{}{
			"uid":    userNode.UserUID,
			"status": userNode.Status,
		})

		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return &domain.UserNode{
			UserUID: record.Values[0].(string),
			Status:  domain.ProfileStatus(record.Values[1].(string)),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*domain.UserNode), nil
}

func checkIfUserExist(uid string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (n:UserNode { uid: $uid }) RETURN n.uid",
		map[string]interface{}{"uid": uid})

	if result != nil && result.Next() && result.Record().Values[0] == uid {
		return true
	}
	return false
}
