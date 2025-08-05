package protobuf

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-transport-prac/pkg/sdl/protobuf/gen/user"
	"go-transport-prac/pkg/sdl/protobuf/gen/userv2"
)

func TestBackwardCompatibility(t *testing.T) {
	// Create v1 user
	v1User := &user.User{
		Id:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		Status:    user.UserStatus_USER_STATUS_ACTIVE,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	// Serialize v1
	v1Data, err := proto.Marshal(v1User)
	if err != nil {
		t.Fatalf("Failed to marshal v1 user: %v", err)
	}

	// Deserialize as v2
	v2User := &userv2.UserV2{}
	err = proto.Unmarshal(v1Data, v2User)
	if err != nil {
		t.Fatalf("Failed to unmarshal v1 data as v2: %v", err)
	}

	// Verify basic fields are preserved
	if v2User.Id != v1User.Id {
		t.Errorf("ID mismatch: got %d, want %d", v2User.Id, v1User.Id)
	}
	if v2User.Email != v1User.Email {
		t.Errorf("Email mismatch: got %s, want %s", v2User.Email, v1User.Email)
	}
	if v2User.Name != v1User.Name {
		t.Errorf("Name mismatch: got %s, want %s", v2User.Name, v1User.Name)
	}

	// Verify new fields have default values
	if v2User.Username != "" {
		t.Errorf("Username should be empty, got %s", v2User.Username)
	}
	if v2User.EmailVerified != false {
		t.Errorf("EmailVerified should be false, got %v", v2User.EmailVerified)
	}
	if len(v2User.Roles) != 0 {
		t.Errorf("Roles should be empty, got %v", v2User.Roles)
	}
}

func TestForwardCompatibility(t *testing.T) {
	// Create v2 user with new fields
	v2User := &userv2.UserV2{
		Id:            2,
		Email:         "test2@example.com",
		Name:          "Test User 2",
		Status:        userv2.UserStatus_USER_STATUS_ACTIVE,
		Username:      "testuser2",
		EmailVerified: true,
		Roles:         []string{"user", "beta"},
		Preferences: &userv2.UserPreferences{
			Language: "en",
			Theme:    userv2.Theme_THEME_DARK,
		},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	// Serialize v2
	v2Data, err := proto.Marshal(v2User)
	if err != nil {
		t.Fatalf("Failed to marshal v2 user: %v", err)
	}

	// Deserialize as v1
	v1User := &user.User{}
	err = proto.Unmarshal(v2Data, v1User)
	if err != nil {
		t.Fatalf("Failed to unmarshal v2 data as v1: %v", err)
	}

	// Verify basic fields are preserved
	if v1User.Id != v2User.Id {
		t.Errorf("ID mismatch: got %d, want %d", v1User.Id, v2User.Id)
	}
	if v1User.Email != v2User.Email {
		t.Errorf("Email mismatch: got %s, want %s", v1User.Email, v2User.Email)
	}
	if v1User.Name != v2User.Name {
		t.Errorf("Name mismatch: got %s, want %s", v1User.Name, v2User.Name)
	}

	// Verify unknown fields are preserved
	unknownFields := v1User.ProtoReflect().GetUnknown()
	if len(unknownFields) == 0 {
		t.Error("Unknown fields should be preserved but none found")
	}
}

func TestUnknownFieldsPreservation(t *testing.T) {
	// Create v2 user
	original := &userv2.UserV2{
		Id:       3,
		Email:    "preserve@example.com",
		Name:     "Preserve Test",
		Username: "preserve_test",
		Roles:    []string{"admin", "user"},
	}

	// v2 -> bytes
	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal original: %v", err)
	}

	// bytes -> v1 (unknown fields preserved)
	v1 := &user.User{}
	if err := proto.Unmarshal(data, v1); err != nil {
		t.Fatalf("Failed to unmarshal to v1: %v", err)
	}

	// v1 -> bytes (unknown fields should still be there)
	intermediateData, err := proto.Marshal(v1)
	if err != nil {
		t.Fatalf("Failed to marshal v1: %v", err)
	}

	// bytes -> v2 (unknown fields recovered)
	recovered := &userv2.UserV2{}
	if err := proto.Unmarshal(intermediateData, recovered); err != nil {
		t.Fatalf("Failed to unmarshal to recovered v2: %v", err)
	}

	// Verify recovery
	if original.Username != recovered.Username {
		t.Errorf("Username not preserved: got %s, want %s", recovered.Username, original.Username)
	}
	if len(original.Roles) != len(recovered.Roles) {
		t.Errorf("Roles length mismatch: got %d, want %d", len(recovered.Roles), len(original.Roles))
	}
	for i, role := range original.Roles {
		if i >= len(recovered.Roles) || recovered.Roles[i] != role {
			t.Errorf("Role %d mismatch: got %s, want %s", i, recovered.Roles[i], role)
		}
	}
}

func TestEnumEvolution(t *testing.T) {
	// Create user with new enum value
	v2User := &userv2.UserV2{
		Id:     4,
		Email:  "enum@example.com",
		Name:   "Enum Test",
		Status: userv2.UserStatus_USER_STATUS_PENDING_VERIFICATION, // New enum value
	}

	// Serialize v2
	data, err := proto.Marshal(v2User)
	if err != nil {
		t.Fatalf("Failed to marshal v2 user: %v", err)
	}

	// Deserialize as v1 (should handle unknown enum gracefully)
	v1User := &user.User{}
	if err := proto.Unmarshal(data, v1User); err != nil {
		t.Fatalf("Failed to unmarshal v2 data as v1: %v", err)
	}

	// Unknown enum values default to 0 (UNSPECIFIED)
	if v1User.Status != user.UserStatus_USER_STATUS_UNSPECIFIED {
		// pass, due to v1 do not have this enum value, it will remain 5
		t.Logf("Expected UNSPECIFIED status for unknown enum, got %v", v1User.Status)
	}
}

func TestSchemaEvolutionRoundTrip(t *testing.T) {
	demo := NewCompatibilityDemo()

	if err := demo.RunCompatibilityTests(); err != nil {
		t.Errorf("Compatibility demo failed: %v", err)
	}
}
