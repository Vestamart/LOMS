package localErr

import "errors"

var ItemNotEnoughErr = errors.New("item not enough")

var SKUNotExistErr = errors.New("sku not exist")

var OrderNotFoundErr = errors.New("order not found")
