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
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(baseUserUuid, tx) && checkIfUserExist(connectUserUuid, tx) {

			records, err := tx.Run("MATCH (n:UserNode { uid: $uid}) RETURN n.status", map[string]interface{}{
				"uid": baseUserUuid,
			})
			connectionStatus := ""
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
				connectionStatus = "REQUEST_SENT"
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
				connectionStatus = "CONNECTED"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $connect_user_uuid  MATCH (u2:UserNode) WHERE u2.uid = $base_user_uuid CREATE (u1)-[r2:CONNECTION {status: $status, date: $date}]->(u2) CREATE (u2)-[r1:CONNECTION {status: $status, date: $date}]->(u1)", map[string]interface{}{
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
				ConnectionResponse: connectionStatus,
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

func (store *ConnectionDBStore) AcceptConnection(requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(requestSenderUser, tx) && checkIfUserExist(requestApprovalUser, tx) {

			records, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $requestSender MATCH (u2:UserNode) WHERE u2.uid = $requestGet  MATCH (u1)-[r1:CONNECTION]->(u2) return r1.status", map[string]interface{}{
				"requestSender": requestSenderUser,
				"requestGet":    requestApprovalUser,
			})
			connectionStatus := ""
			if err != nil {
				return nil, err
			}
			record, err := records.Single()
			if err != nil {
				return nil, err
			}

			status := record.Values[0].(string)
			dateNow := time.Now().Local().Unix()
			if status == "REQUEST_SENT" {
				connectionStatus = "CONNECTED"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $requestSender MATCH (u2:UserNode) WHERE u2.uid = $requestGet match  (u1)-[r1:CONNECTION {status: $status0 }]->(u2) set r1.status = $status1 , r1.date = $date  CREATE  (u2)-[r2:CONNECTION {status: $status1, date: $date}]->(u1)", map[string]interface{}{
					"requestSender": requestSenderUser,
					"requestGet":    requestApprovalUser,
					"status0":       "REQUEST_SENT",
					"status1":       "CONNECTED",
					"date":          dateNow,
				})

				if err != nil {
					return nil, err
				}

			} else if status == "CONNECTED" {
				connectionStatus = "CONNECTION EXISTS"
				if err != nil {
					return nil, err
				}
			}

			// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
			return &pb.NewConnectionResponse{
				BaseUserUUID:       requestSenderUser,
				ConnectUserUUID:    requestApprovalUser,
				ConnectionResponse: connectionStatus,
			}, nil
		} else {
			return &pb.NewConnectionResponse{
				BaseUserUUID:       "",
				ConnectUserUUID:    "",
				ConnectionResponse: "Connection refused - user not found",
			}, nil
		}

	})

	if err != nil {
		return nil, err
	}

	return result.(*pb.NewConnectionResponse), nil
}

func (store *ConnectionDBStore) GetAllConnectionForUser(userUid string) (userNodes []*domain.UserNode, error1 error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if !checkIfUserExist(userUid, tx) {
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		records, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $userUid MATCH (u2:UserNode) WHERE not u2.uid = $userUid match (u2)-[r1:CONNECTION {status: $status}]->(u1) match (u1)-[r2:CONNECTION {status: $status }]->(u2) return u2.uid, u2.status", map[string]interface{}{
			"userUid": userUid,
			"status":  "CONNECTED",
		})

		for records.Next() {
			node := domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Private}
			userNodes = append(userNodes, &node)
		}

		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		return userNodes, nil
	})
	if err != nil {
		return nil, err
	}
	return userNodes, nil
}

func (store *ConnectionDBStore) Register(userNode *domain.UserNode) (*domain.UserNode, error) {
	fmt.Println("[ConnectionDBStore Register]")
	fmt.Println(userNode)
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	
	fmt.Println(session)
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		fmt.Println("linija5")
		if checkIfUserExist(userNode.UserUID, tx) {
			fmt.Println("linija6")
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		fmt.Println("[ConnectionDBStore Register1]")
		fmt.Println(userNode)
		records, err := tx.Run("CREATE (n:UserNode { uid: $uid, status: $status  }) RETURN n.uid, n.status", map[string]interface{}{
			"uid":    userNode.UserUID,
			"status": userNode.Status,
		})
		fmt.Println("TU SAM")
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
