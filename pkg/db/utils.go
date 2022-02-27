package db

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

func attributevalueHasLength(attr types.AttributeValue) bool {
	switch attr := attr.(type) {
	case *types.AttributeValueMemberB:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberBS:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberL:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberM:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberN:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberNS:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberS:
		return len(attr.Value) > 0
	case *types.AttributeValueMemberSS:
		return len(attr.Value) > 0
	default:
		// AttributeValueMemberBOOL
		// AttributeValueMemberNULL
		return true
	}
}

func FilterDynamodbAttributevalueList(attrs []types.AttributeValue) []types.AttributeValue {
	filtered := make([]types.AttributeValue, len(attrs))

	for _, v := range attrs {
		switch v := v.(type) {
		case *types.AttributeValueMemberL:
			v.Value = FilterDynamodbAttributevalueList(v.Value)
			if len(v.Value) > 0 {
				filtered = append(filtered, v)
			}
		case *types.AttributeValueMemberM:
			v.Value = FilterDynamodbAttributevalueMap(v.Value)
			if len(v.Value) > 0 {
				filtered = append(filtered, v)
			}
		default:
			if attributevalueHasLength(v) {
				filtered = append(filtered, v)
			}
		}
	}

	return filtered
}

func FilterDynamodbAttributevalueMap(attrs map[string]types.AttributeValue) map[string]types.AttributeValue {
	filtered := make(map[string]types.AttributeValue, len(attrs))

	for k, v := range attrs {
		switch v := v.(type) {
		case *types.AttributeValueMemberL:
			v.Value = FilterDynamodbAttributevalueList(v.Value)
			if len(v.Value) > 0 {
				filtered[k] = v
			}
		case *types.AttributeValueMemberM:
			v.Value = FilterDynamodbAttributevalueMap(v.Value)
			if len(v.Value) > 0 {
				filtered[k] = v
			}
		default:
			if attributevalueHasLength(v) {
				filtered[k] = v
			}
		}
	}

	return filtered
}
