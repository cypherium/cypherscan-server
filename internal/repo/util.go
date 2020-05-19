package repo

// type createQueryInputOptions struct {
// 	Projection             []string
// 	KeyConditionExpression *expression
// 	IndexName              indexName
// 	Limit                  int
// 	ScanIndexForward       bool
// 	After                  string
// 	Before                 string
// }
// type expression struct {
// 	Fields     []string
// 	Expression string
// 	Values     map[string]*dynamodb.AttributeValue
// }

// // query is a generic function to query
// func (r *Repo) queryOne(option *createQueryInputOptions) (*BlockItem, error) {
// 	option.Limit = 1
// 	queryInput, err := createQueryInput(option)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req, resp := r.Db.QueryRequest(queryInput)
// 	err = req.Send()
// 	if err != nil {
// 		bizutil.HandleError(err, "query failed")
// 		return nil, err
// 	}
// 	if len(resp.Items) < 1 {
// 		return nil, nil
// 	}
// 	item := resp.Items[0]
// 	blockItem := new(BlockItem)
// 	err = dynamodbattribute.UnmarshalMap(item, blockItem)
// 	if err != nil {
// 		return nil, util.NewError(err, "failed to unmarshal")
// 	}
// 	return blockItem, nil
// }

// QueryResult is the struct hold the query result
type QueryResult struct {
	Cursor
	Items []Transaction
}

// CursorPaginationRequest is the generic request for the pagination request
type CursorPaginationRequest struct {
	Before   string
	After    string
	PageSize int
}

// PaginationType is cursor or offset, if offset need to provide latestNo and pageNo, if cursor need to provide before and after
type PaginationType string

const (
	// OffsetType need to provide pageno and latestNo
	OffsetType PaginationType = "Offset"
	// CursorType need to provide before or after
	CursorType PaginationType = "Cursor"
)

// PaginationRequest is the generic request for pagination request
type PaginationRequest struct {
	PageSize int
	Type     PaginationType

	PageNo   int64
	LatestNo int64

	Before int64
	After  int64
}

// const (
// 	// MaxPageSize is the max page size
// 	MaxPageSize = 100
// )

// func normaliseOffsetPaginationRequest(option *PaginationRequest) {
// 	if option.Type == "" {
// 		option.Type = OffsetType
// 	}
// 	if option.PageSize <= 0 || MaxPageSize > 100 {
// 		option.PageSize = MaxPageSize
// 	}
// 	if option.PageNo == 0 {
// 		option.PageNo = 1
// 	}
// 	if option.LatestNo == 0 {
// 		option.LatestNo = util.MaxInt64
// 		option.PageNo = 1
// 	}
// 	if option.After > 0 && option.Before > 0 {
// 		option.Before = 0
// 	}
// 	if option.After == 0 && option.Before == 0 {
// 		option.After = util.MaxInt64
// 	}
// }

// Cursor is the indication of the last of the first one
type Cursor struct {
	First string
	Last  string
}

// type cursor = map[string]*dynamodb.AttributeValue

// // query is a generic function to query
// func (r *Repo) query(option *createQueryInputOptions) (*QueryResult, error) {
// 	queryInput, err := createQueryInput(option)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req, resp := r.Db.QueryRequest(queryInput)
// 	err = req.Send()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if option.Before != "" {
// 		reverse(resp.Items)
// 	}
// 	items := convertDynamodbResultToBlockItems(resp)
// 	cursor, err := getCursor(option, resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &QueryResult{Cursor: *cursor, Items: items}, nil
// }

// func reverse(items []map[string]*dynamodb.AttributeValue) {
// 	i := 0
// 	j := len(items) - 1
// 	for i < j {
// 		items[i], items[j] = items[j], items[i]
// 		i++
// 		j--
// 	}
// }

// func convertDynamodbResultToBlockItems(dynamodbOutput *dynamodb.QueryOutput) []*BlockItem {
// 	ret := make([]*BlockItem, 0, len(dynamodbOutput.Items))
// 	for _, item := range dynamodbOutput.Items {
// 		blockItem := new(BlockItem)
// 		err := dynamodbattribute.UnmarshalMap(item, blockItem)
// 		if err != nil {
// 			bizutil.HandleError(err, "failed to unmarshal")
// 			continue
// 		}
// 		ret = append(ret, blockItem)
// 	}
// 	return ret
// }

// func getCursor(option *createQueryInputOptions, output *dynamodb.QueryOutput) (*Cursor, error) {
// 	if len(output.Items) <= 0 {
// 		return &Cursor{}, nil
// 	}
// 	names := getCursorPropNames(option)
// 	last, errLast := getCursorString(names, output.Items[len(output.Items)-1])
// 	first, errFirst := getCursorString(names, output.Items[0])
// 	if option.Limit > len(output.Items) {
// 		if option.Before != "" {
// 			first = ""
// 			errFirst = nil
// 		} else {
// 			last = ""
// 			errLast = nil
// 		}
// 	}
// 	if errLast != nil || errFirst != nil {
// 		return nil, util.PickFirstFromErrs(errLast, errFirst)
// 	}
// 	return &Cursor{Last: last, First: first}, nil
// }

// func getCursorString(names []string, item map[string]*dynamodb.AttributeValue) (string, error) {
// 	cursor := make(map[string]*dynamodb.AttributeValue)
// 	for _, n := range names {
// 		cursor[n] = item[n]
// 	}

// 	return encodeMapAtrributeValues(cursor)
// }

// func createQueryInput(option *createQueryInputOptions) (*dynamodb.QueryInput, error) {
// 	projectionExpression := getProjectionExpression(option)
// 	indexName := convertIndexNameToAwsString(option.IndexName)
// 	expressionNames := getExpressionNames(option)
// 	exclusiveStartKey, err := getExclusiveStartKey(option)
// 	if err != nil {
// 		return nil, util.NewError(err, "get query input failed")
// 	}
// 	scanIndexForward := getScanIndexForward(option)
// 	return &dynamodb.QueryInput{
// 		KeyConditionExpression:    aws.String(option.KeyConditionExpression.Expression),
// 		ConsistentRead:            aws.Bool(false),
// 		IndexName:                 indexName,
// 		ExpressionAttributeNames:  expressionNames,
// 		ExpressionAttributeValues: option.KeyConditionExpression.Values,
// 		Limit:                aws.Int64(int64(option.Limit)),
// 		ScanIndexForward:     scanIndexForward,
// 		ProjectionExpression: projectionExpression,
// 		TableName:            aws.String(getConst().tableName),
// 		ExclusiveStartKey:    exclusiveStartKey,
// 	}, nil
// }
// func getExclusiveStartKey(option *createQueryInputOptions) (map[string]*dynamodb.AttributeValue, error) {
// 	last, lastErr := decodeMapAttributeValues(option.After)
// 	before, beforeErr := decodeMapAttributeValues(option.Before)
// 	err := util.PickFirstFromErrs(lastErr, beforeErr)
// 	if err != nil {
// 		return nil, util.NewError(err, "Failed to get exclusiveStartKey")
// 	}
// 	exclusiveStartKey := last
// 	if last == nil {
// 		exclusiveStartKey = before
// 	}
// 	return exclusiveStartKey, nil
// }
// func getScanIndexForward(option *createQueryInputOptions) *bool {
// 	scanIndexForward := aws.Bool(option.ScanIndexForward)
// 	if option.After == "" && option.Before != "" {
// 		scanIndexForward = aws.Bool(!option.ScanIndexForward)
// 	}
// 	return scanIndexForward
// }
// func convertIndexNameToAwsString(index indexName) *string {
// 	indexNameStr := string(index)
// 	indexName := &indexNameStr
// 	if *indexName == "" {
// 		indexName = nil
// 	}
// 	return indexName
// }
// func getCursorPropNames(option *createQueryInputOptions) []string {
// 	names := getConst().tablePKFileds
// 	if string(option.IndexName) != "" {
// 		names = uniqStrings(names, getConst().indexesFields[option.IndexName])
// 	}
// 	return names
// }
// func getProjectionExpression(option *createQueryInputOptions) *string {
// 	names := uniqStrings(option.Projection, getCursorPropNames(option))
// 	project := make([]string, 0, len(option.Projection))
// 	for _, p := range names {
// 		project = append(project, fmt.Sprintf("#%s", p))
// 	}
// 	return aws.String(strings.Join(project, ","))
// }

// func getExpressionNames(option *createQueryInputOptions) map[string]*string {
// 	names := uniqStrings(option.Projection, option.KeyConditionExpression.Fields, getCursorPropNames(option))
// 	expressionAttributeNames := make(map[string]*string)
// 	for _, n := range names {
// 		expressionAttributeNames[fmt.Sprintf("#%s", n)] = aws.String(n)
// 	}
// 	return expressionAttributeNames
// }

// func uniqStrings(arrayStrs ...[]string) []string {
// 	smap := make(map[string]struct{})
// 	ret := make([]string, 0)
// 	for _, strs := range arrayStrs {
// 		for _, s := range strs {
// 			smap[s] = struct{}{}
// 		}
// 	}
// 	for k := range smap {
// 		ret = append(ret, k)
// 	}
// 	return ret
// }

// func encodeMapAtrributeValues(attrs map[string]*dynamodb.AttributeValue) (string, error) {
// 	json, err := json.Marshal(attrs)
// 	if err != nil {
// 		return "", util.NewError(err, "Cannot base64 encode attrs: %v", attrs)
// 	}
// 	return base64.StdEncoding.EncodeToString([]byte(json)), nil
// }

// func decodeMapAttributeValues(v string) (map[string]*dynamodb.AttributeValue, error) {
// 	if v == "" {
// 		return nil, nil
// 	}
// 	jsonStr, err := base64.StdEncoding.DecodeString(v)
// 	if err != nil {
// 		return nil, util.NewError(err, "Cannot base64 decode attrs: %s", v)
// 	}
// 	var ret map[string]*dynamodb.AttributeValue
// 	err = json.Unmarshal(jsonStr, &ret)
// 	if err != nil {
// 		return nil, util.NewError(err, "Cannot base64 decode attrs: %s", v)
// 	}
// 	return ret, nil
// }
