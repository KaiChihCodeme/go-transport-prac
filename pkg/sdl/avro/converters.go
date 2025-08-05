package avro

import (
	"fmt"
	"time"
)

// userToAvroMap converts a User struct to an Avro-compatible map
func (m *Manager) userToAvroMap(user User) map[string]interface{} {
	data := map[string]interface{}{
		"id":        user.ID,
		"email":     user.Email,
		"name":      user.Name,
		"status":    string(user.Status),
		"createdAt": user.CreatedAt.UnixMilli(),
		"updatedAt": user.UpdatedAt.UnixMilli(),
	}

	// Handle profile (optional)
	if user.Profile != nil {
		profileData := map[string]interface{}{
			"firstName": user.Profile.FirstName,
			"lastName":  user.Profile.LastName,
			"interests": user.Profile.Interests,
			"metadata":  user.Profile.Metadata,
		}

		// Handle optional phone
		if user.Profile.Phone != nil {
			profileData["phone"] = map[string]interface{}{"string": *user.Profile.Phone}
		} else {
			profileData["phone"] = nil
		}

		// Handle optional address
		if user.Profile.Address != nil {
			addressData := map[string]interface{}{
				"street":     user.Profile.Address.Street,
				"city":       user.Profile.Address.City,
				"state":      user.Profile.Address.State,
				"postalCode": user.Profile.Address.PostalCode,
				"country":    user.Profile.Address.Country,
			}
			profileData["address"] = map[string]interface{}{"com.example.avro.Address": addressData}
		} else {
			profileData["address"] = nil
		}

		data["profile"] = map[string]interface{}{"com.example.avro.Profile": profileData}
	} else {
		data["profile"] = nil
	}

	return data
}

// avroMapToUser converts an Avro map to a User struct
func (m *Manager) avroMapToUser(data map[string]interface{}) (User, error) {
	user := User{
		ID:    toInt64(data["id"]),
		Email: data["email"].(string),
		Name:  data["name"].(string),
		Status: UserStatus(data["status"].(string)),
	}

	// Handle timestamps
	if createdAtMs := data["createdAt"]; createdAtMs != nil {
		user.CreatedAt = time.UnixMilli(toInt64(createdAtMs))
	}
	if updatedAtMs := data["updatedAt"]; updatedAtMs != nil {
		user.UpdatedAt = time.UnixMilli(toInt64(updatedAtMs))
	}

	// Handle profile (optional)
	if profileData := data["profile"]; profileData != nil {
		if profileMap, ok := profileData.(map[string]interface{}); ok {
			if profileValue, exists := profileMap["com.example.avro.Profile"]; exists {
				if profileValueMap, ok := profileValue.(map[string]interface{}); ok {
					profile := &Profile{
						FirstName: profileValueMap["firstName"].(string),
						LastName:  profileValueMap["lastName"].(string),
						Interests: stringSliceFromInterface(profileValueMap["interests"]),
						Metadata:  stringMapFromInterface(profileValueMap["metadata"]),
					}

					// Handle optional phone
					if phoneData := profileValueMap["phone"]; phoneData != nil {
						// Handle different possible formats for union types
						if phoneMap, ok := phoneData.(map[string]interface{}); ok {
							if phoneValue, exists := phoneMap["string"]; exists {
								phoneStr := phoneValue.(string)
								profile.Phone = &phoneStr
							}
						} else if phoneStr, ok := phoneData.(string); ok {
							// Sometimes unions are returned as direct values
							profile.Phone = &phoneStr
						}
					}

					// Handle optional address
					if addressData := profileValueMap["address"]; addressData != nil {
						if addressMap, ok := addressData.(map[string]interface{}); ok {
							if addressValue, exists := addressMap["com.example.avro.Address"]; exists {
								if addressValueMap, ok := addressValue.(map[string]interface{}); ok {
									profile.Address = &Address{
										Street:     addressValueMap["street"].(string),
										City:       addressValueMap["city"].(string),
										State:      addressValueMap["state"].(string),
										PostalCode: addressValueMap["postalCode"].(string),
										Country:    addressValueMap["country"].(string),
									}
								}
							}
						}
					}

					user.Profile = profile
				}
			}
		}
	}

	return user, nil
}

// productToAvroMap converts a Product struct to an Avro-compatible map
func (m *Manager) productToAvroMap(product Product) map[string]interface{} {
	// Price data
	priceData := map[string]interface{}{
		"currency":    product.Price.Currency,
		"amountCents": product.Price.AmountCents,
	}
	if product.Price.DiscountPercentage != nil {
		priceData["discountPercentage"] = map[string]interface{}{"float": *product.Price.DiscountPercentage}
	} else {
		priceData["discountPercentage"] = nil
	}

	// Inventory data
	inventoryData := map[string]interface{}{
		"quantity":       product.Inventory.Quantity,
		"reserved":       product.Inventory.Reserved,
		"available":      product.Inventory.Available,
		"trackInventory": product.Inventory.TrackInventory,
		"reorderLevel":   product.Inventory.ReorderLevel,
		"maxStock":       product.Inventory.MaxStock,
	}

	return map[string]interface{}{
		"id":            product.ID,
		"name":          product.Name,
		"description":   product.Description,
		"sku":           product.SKU,
		"price":         priceData,
		"inventory":     inventoryData,
		"categories":    product.Categories,
		"tags":          product.Tags,
		"status":        string(product.Status),
		"specifications": product.Specifications,
		"createdAt":     product.CreatedAt.UnixMilli(),
		"updatedAt":     product.UpdatedAt.UnixMilli(),
	}
}

// avroMapToProduct converts an Avro map to a Product struct
func (m *Manager) avroMapToProduct(data map[string]interface{}) (Product, error) {
	product := Product{
		ID:          toInt64(data["id"]),
		Name:        data["name"].(string),
		Description: data["description"].(string),
		SKU:         data["sku"].(string),
		Categories:  stringSliceFromInterface(data["categories"]),
		Tags:        stringSliceFromInterface(data["tags"]),
		Status:      ProductStatus(data["status"].(string)),
		Specifications: stringMapFromInterface(data["specifications"]),
	}

	// Handle timestamps  
	if createdAtMs := data["createdAt"]; createdAtMs != nil {
		product.CreatedAt = time.UnixMilli(toInt64(createdAtMs))
	}
	if updatedAtMs := data["updatedAt"]; updatedAtMs != nil {
		product.UpdatedAt = time.UnixMilli(toInt64(updatedAtMs))
	}

	// Handle price
	if priceData, ok := data["price"].(map[string]interface{}); ok {
		product.Price = Price{
			Currency:    priceData["currency"].(string),
			AmountCents: toInt64(priceData["amountCents"]),
		}

		// Handle optional discount
		if discountData := priceData["discountPercentage"]; discountData != nil {
			if discountMap, ok := discountData.(map[string]interface{}); ok {
				if discountValue, exists := discountMap["float"]; exists {
					discount := float32(discountValue.(float64))
					product.Price.DiscountPercentage = &discount
				}
			} else if discountValue, ok := discountData.(float64); ok {
				// Sometimes unions are returned as direct values
				discount := float32(discountValue)
				product.Price.DiscountPercentage = &discount
			} else if discountValue, ok := discountData.(float32); ok {
				product.Price.DiscountPercentage = &discountValue
			}
		}
	}

	// Handle inventory
	if inventoryData, ok := data["inventory"].(map[string]interface{}); ok {
		product.Inventory = Inventory{
			Quantity:       toInt32(inventoryData["quantity"]),
			Reserved:       toInt32(inventoryData["reserved"]),
			Available:      toInt32(inventoryData["available"]),
			TrackInventory: inventoryData["trackInventory"].(bool),
			ReorderLevel:   toInt32(inventoryData["reorderLevel"]),
			MaxStock:       toInt32(inventoryData["maxStock"]),
		}
	}

	return product, nil
}

// Helper functions

// toInt64 safely converts various numeric types to int64
func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	default:
		return 0
	}
}

// toInt32 safely converts various numeric types to int32
func toInt32(v interface{}) int32 {
	switch val := v.(type) {
	case int:
		return int32(val)
	case int32:
		return val
	case int64:
		return int32(val)
	case float64:
		return int32(val)
	default:
		return 0
	}
}

// stringSliceFromInterface converts an interface{} to []string
func stringSliceFromInterface(data interface{}) []string {
	if data == nil {
		return []string{}
	}

	if slice, ok := data.([]interface{}); ok {
		result := make([]string, len(slice))
		for i, item := range slice {
			result[i] = item.(string)
		}
		return result
	}

	return []string{}
}

// stringMapFromInterface converts an interface{} to map[string]string
func stringMapFromInterface(data interface{}) map[string]string {
	if data == nil {
		return map[string]string{}
	}

	if m, ok := data.(map[string]interface{}); ok {
		result := make(map[string]string)
		for k, v := range m {
			result[k] = v.(string)
		}
		return result
	}

	return map[string]string{}
}

// CompareData compares two interface{} values for testing
func CompareData(a, b interface{}) error {
	// This is a simplified comparison - in production you'd want more robust comparison
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	
	if aStr != bStr {
		return fmt.Errorf("data mismatch: %s != %s", aStr, bStr)
	}
	
	return nil
}