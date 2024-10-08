// Code generated by ObjectBox; DO NOT EDIT.
// Learn more about defining entities and generating this file - visit https://golang.objectbox.io/entity-annotations

package main

import (
	"errors"
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type user_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var userBinding = user_EntityInfo{
	Entity: objectbox.Entity{
		Id: 1,
	},
	Uid: 2723288463117073965,
}

// user_ contains type-based Property helpers to facilitate some common operations such as Queries.
var user_ = struct {
	id           *objectbox.PropertyUint64
	username     *objectbox.PropertyString
	passwordHash *objectbox.PropertyByteVector
}{
	id: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &userBinding.Entity,
		},
	},
	username: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &userBinding.Entity,
		},
	},
	passwordHash: &objectbox.PropertyByteVector{
		BaseProperty: &objectbox.BaseProperty{
			Id:     3,
			Entity: &userBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (user_EntityInfo) GeneratorVersion() int {
	return 6
}

// AddToModel is called by ObjectBox during model build
func (user_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("user", 1, 2723288463117073965)
	model.Property("id", 6, 1, 8209159007819727508)
	model.PropertyFlags(1)
	model.Property("username", 9, 2, 8481836100675543267)
	model.Property("passwordHash", 23, 3, 1231099847571825165)
	model.EntityLastPropertyId(3, 1231099847571825165)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (user_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*user).id, nil
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (user_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*user).id = id
	return nil
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (user_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (user_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	obj := object.(*user)
	var offsetusername = fbutils.CreateStringOffset(fbb, obj.username)
	var offsetpasswordHash = fbutils.CreateByteVectorOffset(fbb, obj.passwordHash)

	// build the FlatBuffers object
	fbb.StartObject(3)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetUOffsetTSlot(fbb, 1, offsetusername)
	fbutils.SetUOffsetTSlot(fbb, 2, offsetpasswordHash)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (user_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	if len(bytes) == 0 { // sanity check, should "never" happen
		return nil, errors.New("can't deserialize an object of type 'user' - no data received")
	}

	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	var propid = table.GetUint64Slot(4, 0)

	return &user{
		id:           propid,
		username:     fbutils.GetStringSlot(table, 6),
		passwordHash: fbutils.GetByteVectorSlot(table, 8),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (user_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*user, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (user_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	if object == nil {
		return append(slice.([]*user), nil)
	}
	return append(slice.([]*user), object.(*user))
}

// Box provides CRUD access to user objects
type userBox struct {
	*objectbox.Box
}

// BoxForuser opens a box of user objects
func BoxForuser(ob *objectbox.ObjectBox) *userBox {
	return &userBox{
		Box: ob.InternalBox(1),
	}
}

// Put synchronously inserts/updates a single object.
// In case the id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the user.id property on the passed object will be assigned the new ID as well.
func (box *userBox) Put(object *user) (uint64, error) {
	return box.Box.Put(object)
}

// Insert synchronously inserts a single object. As opposed to Put, Insert will fail if given an ID that already exists.
// In case the id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the user.id property on the passed object will be assigned the new ID as well.
func (box *userBox) Insert(object *user) (uint64, error) {
	return box.Box.Insert(object)
}

// Update synchronously updates a single object.
// As opposed to Put, Update will fail if an object with the same ID is not found in the database.
func (box *userBox) Update(object *user) error {
	return box.Box.Update(object)
}

// PutAsync asynchronously inserts/updates a single object.
// Deprecated: use box.Async().Put() instead
func (box *userBox) PutAsync(object *user) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the user.id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the user.id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *userBox) PutMany(objects []*user) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *userBox) Get(id uint64) (*user, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*user), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is nil
func (box *userBox) GetMany(ids ...uint64) ([]*user, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*user), nil
}

// GetManyExisting reads multiple objects at once, skipping those that do not exist.
func (box *userBox) GetManyExisting(ids ...uint64) ([]*user, error) {
	objects, err := box.Box.GetManyExisting(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*user), nil
}

// GetAll reads all stored objects
func (box *userBox) GetAll() ([]*user, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*user), nil
}

// Remove deletes a single object
func (box *userBox) Remove(object *user) error {
	return box.Box.Remove(object)
}

// RemoveMany deletes multiple objects at once.
// Returns the number of deleted object or error on failure.
// Note that this method will not fail if an object is not found (e.g. already removed).
// In case you need to strictly check whether all of the objects exist before removing them,
// you can execute multiple box.Contains() and box.Remove() inside a single write transaction.
func (box *userBox) RemoveMany(objects ...*user) (uint64, error) {
	var ids = make([]uint64, len(objects))
	for k, object := range objects {
		ids[k] = object.id
	}
	return box.Box.RemoveIds(ids...)
}

// Creates a query with the given conditions. Use the fields of the user_ struct to create conditions.
// Keep the *userQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *userBox) Query(conditions ...objectbox.Condition) *userQuery {
	return &userQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the user_ struct to create conditions.
// Keep the *userQuery if you intend to execute the query multiple times.
func (box *userBox) QueryOrError(conditions ...objectbox.Condition) (*userQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &userQuery{query}, nil
	}
}

// Async provides access to the default Async Box for asynchronous operations. See userAsyncBox for more information.
func (box *userBox) Async() *userAsyncBox {
	return &userAsyncBox{AsyncBox: box.Box.Async()}
}

// userAsyncBox provides asynchronous operations on user objects.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type userAsyncBox struct {
	*objectbox.AsyncBox
}

// AsyncBoxForuser creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use userBox::Async() which takes care of resource management and doesn't require closing.
func AsyncBoxForuser(ob *objectbox.ObjectBox, timeoutMs uint64) *userAsyncBox {
	var async, err = objectbox.NewAsyncBox(ob, 1, timeoutMs)
	if err != nil {
		panic("Could not create async box for entity ID 1: %s" + err.Error())
	}
	return &userAsyncBox{AsyncBox: async}
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the id property on the passed object will be assigned the new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (asyncBox *userAsyncBox) Put(object *user) (uint64, error) {
	return asyncBox.AsyncBox.Put(object)
}

// Insert a single object asynchronously.
// The id property on the passed object will be assigned the new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (asyncBox *userAsyncBox) Insert(object *user) (id uint64, err error) {
	return asyncBox.AsyncBox.Insert(object)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (asyncBox *userAsyncBox) Update(object *user) error {
	return asyncBox.AsyncBox.Update(object)
}

// Remove deletes a single object asynchronously.
func (asyncBox *userAsyncBox) Remove(object *user) error {
	return asyncBox.AsyncBox.Remove(object)
}

// Query provides a way to search stored objects
//
// For example, you can find all user which id is either 42 or 47:
//
// box.Query(user_.id.In(42, 47)).Find()
type userQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *userQuery) Find() ([]*user, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*user), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *userQuery) Offset(offset uint64) *userQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *userQuery) Limit(limit uint64) *userQuery {
	query.Query.Limit(limit)
	return query
}

type product_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var productBinding = product_EntityInfo{
	Entity: objectbox.Entity{
		Id: 2,
	},
	Uid: 758365456359428036,
}

// product_ contains type-based Property helpers to facilitate some common operations such as Queries.
var product_ = struct {
	id        *objectbox.PropertyUint64
	ProductID *objectbox.PropertyInt
	Category  *objectbox.PropertyString
	Name      *objectbox.PropertyString
	ImageURL  *objectbox.PropertyString
}{
	id: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &productBinding.Entity,
		},
	},
	ProductID: &objectbox.PropertyInt{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &productBinding.Entity,
		},
	},
	Category: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     3,
			Entity: &productBinding.Entity,
		},
	},
	Name: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     4,
			Entity: &productBinding.Entity,
		},
	},
	ImageURL: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     5,
			Entity: &productBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (product_EntityInfo) GeneratorVersion() int {
	return 6
}

// AddToModel is called by ObjectBox during model build
func (product_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("product", 2, 758365456359428036)
	model.Property("id", 6, 1, 205342181156815257)
	model.PropertyFlags(1)
	model.Property("ProductID", 6, 2, 2837665624995909300)
	model.Property("Category", 9, 3, 5517858282805562198)
	model.Property("Name", 9, 4, 1335036859670648677)
	model.Property("ImageURL", 9, 5, 5812541690724631292)
	model.EntityLastPropertyId(5, 5812541690724631292)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (product_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*product).id, nil
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (product_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*product).id = id
	return nil
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (product_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (product_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	obj := object.(*product)
	var offsetCategory = fbutils.CreateStringOffset(fbb, obj.Category)
	var offsetName = fbutils.CreateStringOffset(fbb, obj.Name)
	var offsetImageURL = fbutils.CreateStringOffset(fbb, obj.ImageURL)

	// build the FlatBuffers object
	fbb.StartObject(5)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetInt64Slot(fbb, 1, int64(obj.ProductID))
	fbutils.SetUOffsetTSlot(fbb, 2, offsetCategory)
	fbutils.SetUOffsetTSlot(fbb, 3, offsetName)
	fbutils.SetUOffsetTSlot(fbb, 4, offsetImageURL)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (product_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	if len(bytes) == 0 { // sanity check, should "never" happen
		return nil, errors.New("can't deserialize an object of type 'product' - no data received")
	}

	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	var propid = table.GetUint64Slot(4, 0)

	return &product{
		id:        propid,
		ProductID: fbutils.GetIntSlot(table, 6),
		Category:  fbutils.GetStringSlot(table, 8),
		Name:      fbutils.GetStringSlot(table, 10),
		ImageURL:  fbutils.GetStringSlot(table, 12),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (product_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*product, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (product_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	if object == nil {
		return append(slice.([]*product), nil)
	}
	return append(slice.([]*product), object.(*product))
}

// Box provides CRUD access to product objects
type productBox struct {
	*objectbox.Box
}

// BoxForproduct opens a box of product objects
func BoxForproduct(ob *objectbox.ObjectBox) *productBox {
	return &productBox{
		Box: ob.InternalBox(2),
	}
}

// Put synchronously inserts/updates a single object.
// In case the id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the product.id property on the passed object will be assigned the new ID as well.
func (box *productBox) Put(object *product) (uint64, error) {
	return box.Box.Put(object)
}

// Insert synchronously inserts a single object. As opposed to Put, Insert will fail if given an ID that already exists.
// In case the id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the product.id property on the passed object will be assigned the new ID as well.
func (box *productBox) Insert(object *product) (uint64, error) {
	return box.Box.Insert(object)
}

// Update synchronously updates a single object.
// As opposed to Put, Update will fail if an object with the same ID is not found in the database.
func (box *productBox) Update(object *product) error {
	return box.Box.Update(object)
}

// PutAsync asynchronously inserts/updates a single object.
// Deprecated: use box.Async().Put() instead
func (box *productBox) PutAsync(object *product) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the product.id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the product.id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *productBox) PutMany(objects []*product) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *productBox) Get(id uint64) (*product, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*product), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is nil
func (box *productBox) GetMany(ids ...uint64) ([]*product, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*product), nil
}

// GetManyExisting reads multiple objects at once, skipping those that do not exist.
func (box *productBox) GetManyExisting(ids ...uint64) ([]*product, error) {
	objects, err := box.Box.GetManyExisting(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*product), nil
}

// GetAll reads all stored objects
func (box *productBox) GetAll() ([]*product, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*product), nil
}

// Remove deletes a single object
func (box *productBox) Remove(object *product) error {
	return box.Box.Remove(object)
}

// RemoveMany deletes multiple objects at once.
// Returns the number of deleted object or error on failure.
// Note that this method will not fail if an object is not found (e.g. already removed).
// In case you need to strictly check whether all of the objects exist before removing them,
// you can execute multiple box.Contains() and box.Remove() inside a single write transaction.
func (box *productBox) RemoveMany(objects ...*product) (uint64, error) {
	var ids = make([]uint64, len(objects))
	for k, object := range objects {
		ids[k] = object.id
	}
	return box.Box.RemoveIds(ids...)
}

// Creates a query with the given conditions. Use the fields of the product_ struct to create conditions.
// Keep the *productQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *productBox) Query(conditions ...objectbox.Condition) *productQuery {
	return &productQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the product_ struct to create conditions.
// Keep the *productQuery if you intend to execute the query multiple times.
func (box *productBox) QueryOrError(conditions ...objectbox.Condition) (*productQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &productQuery{query}, nil
	}
}

// Async provides access to the default Async Box for asynchronous operations. See productAsyncBox for more information.
func (box *productBox) Async() *productAsyncBox {
	return &productAsyncBox{AsyncBox: box.Box.Async()}
}

// productAsyncBox provides asynchronous operations on product objects.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type productAsyncBox struct {
	*objectbox.AsyncBox
}

// AsyncBoxForproduct creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use productBox::Async() which takes care of resource management and doesn't require closing.
func AsyncBoxForproduct(ob *objectbox.ObjectBox, timeoutMs uint64) *productAsyncBox {
	var async, err = objectbox.NewAsyncBox(ob, 2, timeoutMs)
	if err != nil {
		panic("Could not create async box for entity ID 2: %s" + err.Error())
	}
	return &productAsyncBox{AsyncBox: async}
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the id property on the passed object will be assigned the new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (asyncBox *productAsyncBox) Put(object *product) (uint64, error) {
	return asyncBox.AsyncBox.Put(object)
}

// Insert a single object asynchronously.
// The id property on the passed object will be assigned the new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (asyncBox *productAsyncBox) Insert(object *product) (id uint64, err error) {
	return asyncBox.AsyncBox.Insert(object)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (asyncBox *productAsyncBox) Update(object *product) error {
	return asyncBox.AsyncBox.Update(object)
}

// Remove deletes a single object asynchronously.
func (asyncBox *productAsyncBox) Remove(object *product) error {
	return asyncBox.AsyncBox.Remove(object)
}

// Query provides a way to search stored objects
//
// For example, you can find all product which id is either 42 or 47:
//
// box.Query(product_.id.In(42, 47)).Find()
type productQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *productQuery) Find() ([]*product, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*product), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *productQuery) Offset(offset uint64) *productQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *productQuery) Limit(limit uint64) *productQuery {
	query.Query.Limit(limit)
	return query
}

type store_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var storeBinding = store_EntityInfo{
	Entity: objectbox.Entity{
		Id: 3,
	},
	Uid: 951846619275233843,
}

// store_ contains type-based Property helpers to facilitate some common operations such as Queries.
var store_ = struct {
	ID      *objectbox.PropertyUint64
	Name    *objectbox.PropertyString
	Address *objectbox.PropertyString
	Grid    *objectbox.PropertyByteVector
	Width   *objectbox.PropertyInt
	Start_X *objectbox.PropertyInt
	Start_Y *objectbox.PropertyInt
	Owner   *objectbox.PropertyUint64
}{
	ID: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &storeBinding.Entity,
		},
	},
	Name: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &storeBinding.Entity,
		},
	},
	Address: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     3,
			Entity: &storeBinding.Entity,
		},
	},
	Grid: &objectbox.PropertyByteVector{
		BaseProperty: &objectbox.BaseProperty{
			Id:     4,
			Entity: &storeBinding.Entity,
		},
	},
	Width: &objectbox.PropertyInt{
		BaseProperty: &objectbox.BaseProperty{
			Id:     5,
			Entity: &storeBinding.Entity,
		},
	},
	Start_X: &objectbox.PropertyInt{
		BaseProperty: &objectbox.BaseProperty{
			Id:     6,
			Entity: &storeBinding.Entity,
		},
	},
	Start_Y: &objectbox.PropertyInt{
		BaseProperty: &objectbox.BaseProperty{
			Id:     7,
			Entity: &storeBinding.Entity,
		},
	},
	Owner: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     8,
			Entity: &storeBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (store_EntityInfo) GeneratorVersion() int {
	return 6
}

// AddToModel is called by ObjectBox during model build
func (store_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("store", 3, 951846619275233843)
	model.Property("ID", 6, 1, 10942356598834046)
	model.PropertyFlags(1)
	model.Property("Name", 9, 2, 3751809351370030338)
	model.Property("Address", 9, 3, 2998089498150673737)
	model.Property("Grid", 23, 4, 8570255972576011839)
	model.Property("Width", 6, 5, 4492379083027535731)
	model.Property("Start_X", 6, 6, 4655222207200051383)
	model.Property("Start_Y", 6, 7, 8094724709565254531)
	model.Property("Owner", 6, 8, 5751424598683235360)
	model.PropertyFlags(8192)
	model.EntityLastPropertyId(8, 5751424598683235360)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (store_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*store).ID, nil
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (store_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*store).ID = id
	return nil
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (store_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (store_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	obj := object.(*store)
	var offsetName = fbutils.CreateStringOffset(fbb, obj.Name)
	var offsetAddress = fbutils.CreateStringOffset(fbb, obj.Address)
	var offsetGrid = fbutils.CreateByteVectorOffset(fbb, obj.Grid)

	// build the FlatBuffers object
	fbb.StartObject(8)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetUOffsetTSlot(fbb, 1, offsetName)
	fbutils.SetUOffsetTSlot(fbb, 2, offsetAddress)
	fbutils.SetInt64Slot(fbb, 4, int64(obj.Width))
	fbutils.SetUOffsetTSlot(fbb, 3, offsetGrid)
	fbutils.SetInt64Slot(fbb, 5, int64(obj.Start.X))
	fbutils.SetInt64Slot(fbb, 6, int64(obj.Start.Y))
	fbutils.SetUint64Slot(fbb, 7, obj.Owner)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (store_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	if len(bytes) == 0 { // sanity check, should "never" happen
		return nil, errors.New("can't deserialize an object of type 'store' - no data received")
	}

	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	var propID = table.GetUint64Slot(4, 0)

	return &store{
		ID:      propID,
		Name:    fbutils.GetStringSlot(table, 6),
		Address: fbutils.GetStringSlot(table, 8),
		Width:   fbutils.GetIntSlot(table, 12),
		Grid:    fbutils.GetByteVectorSlot(table, 10),
		Start: point{
			X: fbutils.GetIntSlot(table, 14),
			Y: fbutils.GetIntSlot(table, 16),
		},
		Owner: fbutils.GetUint64Slot(table, 18),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (store_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*store, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (store_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	if object == nil {
		return append(slice.([]*store), nil)
	}
	return append(slice.([]*store), object.(*store))
}

// Box provides CRUD access to store objects
type storeBox struct {
	*objectbox.Box
}

// BoxForstore opens a box of store objects
func BoxForstore(ob *objectbox.ObjectBox) *storeBox {
	return &storeBox{
		Box: ob.InternalBox(3),
	}
}

// Put synchronously inserts/updates a single object.
// In case the ID is not specified, it would be assigned automatically (auto-increment).
// When inserting, the store.ID property on the passed object will be assigned the new ID as well.
func (box *storeBox) Put(object *store) (uint64, error) {
	return box.Box.Put(object)
}

// Insert synchronously inserts a single object. As opposed to Put, Insert will fail if given an ID that already exists.
// In case the ID is not specified, it would be assigned automatically (auto-increment).
// When inserting, the store.ID property on the passed object will be assigned the new ID as well.
func (box *storeBox) Insert(object *store) (uint64, error) {
	return box.Box.Insert(object)
}

// Update synchronously updates a single object.
// As opposed to Put, Update will fail if an object with the same ID is not found in the database.
func (box *storeBox) Update(object *store) error {
	return box.Box.Update(object)
}

// PutAsync asynchronously inserts/updates a single object.
// Deprecated: use box.Async().Put() instead
func (box *storeBox) PutAsync(object *store) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case IDs are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the store.ID property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the store.ID assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *storeBox) PutMany(objects []*store) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *storeBox) Get(id uint64) (*store, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*store), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is nil
func (box *storeBox) GetMany(ids ...uint64) ([]*store, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*store), nil
}

// GetManyExisting reads multiple objects at once, skipping those that do not exist.
func (box *storeBox) GetManyExisting(ids ...uint64) ([]*store, error) {
	objects, err := box.Box.GetManyExisting(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*store), nil
}

// GetAll reads all stored objects
func (box *storeBox) GetAll() ([]*store, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*store), nil
}

// Remove deletes a single object
func (box *storeBox) Remove(object *store) error {
	return box.Box.Remove(object)
}

// RemoveMany deletes multiple objects at once.
// Returns the number of deleted object or error on failure.
// Note that this method will not fail if an object is not found (e.g. already removed).
// In case you need to strictly check whether all of the objects exist before removing them,
// you can execute multiple box.Contains() and box.Remove() inside a single write transaction.
func (box *storeBox) RemoveMany(objects ...*store) (uint64, error) {
	var ids = make([]uint64, len(objects))
	for k, object := range objects {
		ids[k] = object.ID
	}
	return box.Box.RemoveIds(ids...)
}

// Creates a query with the given conditions. Use the fields of the store_ struct to create conditions.
// Keep the *storeQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *storeBox) Query(conditions ...objectbox.Condition) *storeQuery {
	return &storeQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the store_ struct to create conditions.
// Keep the *storeQuery if you intend to execute the query multiple times.
func (box *storeBox) QueryOrError(conditions ...objectbox.Condition) (*storeQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &storeQuery{query}, nil
	}
}

// Async provides access to the default Async Box for asynchronous operations. See storeAsyncBox for more information.
func (box *storeBox) Async() *storeAsyncBox {
	return &storeAsyncBox{AsyncBox: box.Box.Async()}
}

// storeAsyncBox provides asynchronous operations on store objects.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type storeAsyncBox struct {
	*objectbox.AsyncBox
}

// AsyncBoxForstore creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use storeBox::Async() which takes care of resource management and doesn't require closing.
func AsyncBoxForstore(ob *objectbox.ObjectBox, timeoutMs uint64) *storeAsyncBox {
	var async, err = objectbox.NewAsyncBox(ob, 3, timeoutMs)
	if err != nil {
		panic("Could not create async box for entity ID 3: %s" + err.Error())
	}
	return &storeAsyncBox{AsyncBox: async}
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the ID property on the passed object will be assigned the new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (asyncBox *storeAsyncBox) Put(object *store) (uint64, error) {
	return asyncBox.AsyncBox.Put(object)
}

// Insert a single object asynchronously.
// The ID property on the passed object will be assigned the new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (asyncBox *storeAsyncBox) Insert(object *store) (id uint64, err error) {
	return asyncBox.AsyncBox.Insert(object)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (asyncBox *storeAsyncBox) Update(object *store) error {
	return asyncBox.AsyncBox.Update(object)
}

// Remove deletes a single object asynchronously.
func (asyncBox *storeAsyncBox) Remove(object *store) error {
	return asyncBox.AsyncBox.Remove(object)
}

// Query provides a way to search stored objects
//
// For example, you can find all store which ID is either 42 or 47:
//
// box.Query(store_.ID.In(42, 47)).Find()
type storeQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *storeQuery) Find() ([]*store, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*store), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *storeQuery) Offset(offset uint64) *storeQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *storeQuery) Limit(limit uint64) *storeQuery {
	query.Query.Limit(limit)
	return query
}
