package protobuf

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-transport-prac/pkg/sdl/protobuf/gen/user"
	"go-transport-prac/pkg/sdl/protobuf/gen/userv2"
)

// CompatibilityDemo demonstrates backward and forward compatibility features of Protocol Buffers
type CompatibilityDemo struct {
	manager *Manager
}

// NewCompatibilityDemo creates a new compatibility demo instance
func NewCompatibilityDemo() *CompatibilityDemo {
	return &CompatibilityDemo{
		manager: NewManager(),
	}
}

// RunCompatibilityTests runs all compatibility tests
func (c *CompatibilityDemo) RunCompatibilityTests() error {
	fmt.Println("=== Protocol Buffers Compatibility Demonstrations ===")

	if err := c.BackwardCompatibilityDemo(); err != nil {
		return fmt.Errorf("backward compatibility demo failed: %w", err)
	}

	if err := c.ForwardCompatibilityDemo(); err != nil {
		return fmt.Errorf("forward compatibility demo failed: %w", err)
	}

	if err := c.UnknownFieldsDemo(); err != nil {
		return fmt.Errorf("unknown fields demo failed: %w", err)
	}

	if err := c.FieldEvolutionDemo(); err != nil {
		return fmt.Errorf("field evolution demo failed: %w", err)
	}

	return nil
}

// BackwardCompatibilityDemo shows how new code can read old data
func (c *CompatibilityDemo) BackwardCompatibilityDemo() error {
	fmt.Println("--- Backward Compatibility Demo ---")
	fmt.Println("Testing: New code (UserV2) reading old data (User)")

	// Create a v1 user (old format)
	oldUser := &user.User{
		Id:    1,
		Email: "alice@example.com",
		Name:  "Alice Smith",
		Status: user.UserStatus_USER_STATUS_ACTIVE,
		Profile: &user.Profile{
			FirstName: "Alice",
			LastName:  "Smith",
			Phone:     "+1-555-0456",
			Address: &user.Address{
				Street:     "456 Oak Ave",
				City:       "Portland",
				State:      "OR",
				PostalCode: "97201",
				Country:    "USA",
			},
			Interests: []string{"photography", "hiking"},
			Metadata: map[string]string{
				"preferred_language": "en",
				"timezone":          "America/Los_Angeles",
			},
		},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	// Serialize with v1 format
	oldData, err := proto.Marshal(oldUser)
	if err != nil {
		return fmt.Errorf("failed to serialize v1 user: %w", err)
	}

	fmt.Printf("V1 User serialized size: %d bytes\n", len(oldData))

	// Deserialize with v2 format (should work with backward compatibility)
	newUser := &userv2.UserV2{}
	if err := proto.Unmarshal(oldData, newUser); err != nil {
		return fmt.Errorf("failed to deserialize v1 data with v2 struct: %w", err)
	}

	fmt.Printf("V2 User deserialized from V1 data:\n")
	fmt.Printf("  ID: %d\n", newUser.Id)
	fmt.Printf("  Email: %s\n", newUser.Email)
	fmt.Printf("  Name: %s\n", newUser.Name)
	fmt.Printf("  Username: %s (new field - empty as expected)\n", newUser.Username)
	fmt.Printf("  EmailVerified: %v (new field - default false)\n", newUser.EmailVerified)
	fmt.Printf("  Roles: %v (new field - empty slice)\n", newUser.Roles)
	fmt.Printf("  Preferences: %v (new field - nil as expected)\n", newUser.Preferences)

	fmt.Println("✓ Backward compatibility test successful")
	return nil
}

// ForwardCompatibilityDemo shows how old code can read new data
func (c *CompatibilityDemo) ForwardCompatibilityDemo() error {
	fmt.Println("--- Forward Compatibility Demo ---")
	fmt.Println("Testing: Old code (User) reading new data (UserV2)")

	// Create a v2 user with new fields
	newUser := &userv2.UserV2{
		Id:    2,
		Email: "bob@example.com",
		Name:  "Bob Johnson",
		Status: userv2.UserStatus_USER_STATUS_ACTIVE,
		Profile: &userv2.Profile{
			FirstName: "Bob",
			LastName:  "Johnson",
			Phone:     "+1-555-0789",
			Address: &userv2.Address{
				Street:     "789 Pine St",
				City:       "Seattle",
				State:      "WA",
				PostalCode: "98101",
				Country:    "USA",
			},
			Interests: []string{"cooking", "gaming"},
			Metadata: map[string]string{
				"preferred_language": "en",
				"timezone":          "America/Los_Angeles",
			},
			// Bio and BirthDate are v2-only fields, not available in Profile
		},
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		Username:      "bob_johnson",
		EmailVerified: true,
		Roles:         []string{"user", "beta_tester"},
		Preferences: &userv2.UserPreferences{
			Language:          "en",
			Timezone:          "America/Los_Angeles",
			MarketingEmails:   false,
			PushNotifications: true,
			Theme:             userv2.Theme_THEME_DARK,
		},
		AvatarUrl: "https://example.com/avatars/bob.jpg",
	}

	// Serialize with v2 format
	newData, err := proto.Marshal(newUser)
	if err != nil {
		return fmt.Errorf("failed to serialize v2 user: %w", err)
	}

	fmt.Printf("V2 User serialized size: %d bytes\n", len(newData))

	// Deserialize with v1 format (should work with forward compatibility)
	oldUser := &user.User{}
	if err := proto.Unmarshal(newData, oldUser); err != nil {
		return fmt.Errorf("failed to deserialize v2 data with v1 struct: %w", err)
	}

	fmt.Printf("V1 User deserialized from V2 data:\n")
	fmt.Printf("  ID: %d\n", oldUser.Id)
	fmt.Printf("  Email: %s\n", oldUser.Email)
	fmt.Printf("  Name: %s\n", oldUser.Name)
	fmt.Printf("  Status: %v\n", oldUser.Status)
	fmt.Printf("  Profile contains standard fields only (v2 bio field not accessible)\n")
	
	// Show that unknown fields are preserved
	unknownFields := oldUser.ProtoReflect().GetUnknown()
	fmt.Printf("  Unknown fields preserved: %d bytes\n", len(unknownFields))

	fmt.Println("✓ Forward compatibility test successful")
	return nil
}

// UnknownFieldsDemo demonstrates how unknown fields are handled
func (c *CompatibilityDemo) UnknownFieldsDemo() error {
	fmt.Println("--- Unknown Fields Preservation Demo ---")

	// Create v2 user with new fields
	v2User := &userv2.UserV2{
		Id:       3,
		Email:    "charlie@example.com", 
		Name:     "Charlie Brown",
		Username: "charlie_b",
		Roles:    []string{"admin"},
	}

	// Serialize v2 -> bytes
	v2Data, err := proto.Marshal(v2User)
	if err != nil {
		return err
	}

	// Deserialize to v1 (loses new fields but preserves them as unknown)
	v1User := &user.User{}
	if err := proto.Unmarshal(v2Data, v1User); err != nil {
		return err
	}

	// Serialize v1 back to bytes (unknown fields should be preserved)
	v1Data, err := proto.Marshal(v1User)
	if err != nil {
		return err
	}

	// Deserialize back to v2 (should recover the unknown fields)
	recoveredV2User := &userv2.UserV2{}
	if err := proto.Unmarshal(v1Data, recoveredV2User); err != nil {
		return err
	}

	fmt.Printf("Original v2 username: %s\n", v2User.Username)
	fmt.Printf("Recovered v2 username: %s\n", recoveredV2User.Username)
	fmt.Printf("Original v2 roles: %v\n", v2User.Roles)
	fmt.Printf("Recovered v2 roles: %v\n", recoveredV2User.Roles)

	if v2User.Username == recoveredV2User.Username && 
	   len(v2User.Roles) == len(recoveredV2User.Roles) &&
	   v2User.Roles[0] == recoveredV2User.Roles[0] {
		fmt.Println("✓ Unknown fields preservation test successful")
	} else {
		return fmt.Errorf("unknown fields were not properly preserved")
	}

	return nil
}

// FieldEvolutionDemo shows different types of schema evolution
func (c *CompatibilityDemo) FieldEvolutionDemo() error {
	fmt.Println("--- Field Evolution Demo ---")

	// Show different evolution patterns
	fmt.Println("Evolution patterns demonstrated:")
	fmt.Println("1. Adding new optional fields (username, email_verified)")
	fmt.Println("2. Adding new repeated fields (roles)")
	fmt.Println("3. Adding new nested messages (preferences)")
	fmt.Println("4. Adding new enum values (USER_STATUS_PENDING_VERIFICATION)")
	fmt.Println("5. Extending existing messages (Profile with bio, birth_date)")
	fmt.Println("6. Reserved field numbers for future use (field 12)")
	
	// Demonstrate enum evolution
	fmt.Println("\nEnum evolution example:")
	
	// Create user with new enum value
	userWithNewStatus := &userv2.UserV2{
		Id:     4,
		Email:  "pending@example.com",
		Name:   "Pending User",
		Status: userv2.UserStatus_USER_STATUS_PENDING_VERIFICATION, // New enum value
	}

	data, err := proto.Marshal(userWithNewStatus)
	if err != nil {
		return err
	}

	// Deserialize with old version (will get UNSPECIFIED for unknown enum)
	oldUserWithUnknownEnum := &user.User{}
	if err := proto.Unmarshal(data, oldUserWithUnknownEnum); err != nil {
		return err
	}

	fmt.Printf("New enum value: %v\n", userWithNewStatus.Status)
	fmt.Printf("Old code sees: %v (safely defaults to UNSPECIFIED)\n", oldUserWithUnknownEnum.Status)

	fmt.Println("✓ Field evolution demonstration completed")
	return nil
}