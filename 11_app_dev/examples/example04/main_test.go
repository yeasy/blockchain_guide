package main

import (
	"testing"

	"github.com/yeasy/blockchain_guide/11_app_dev/internal/testledger"
)

func TestEnrollmentAndDiplomaLifecycleRejectsDuplicates(t *testing.T) {
	ctx, stub := testledger.New("school-1", 100)
	contract := new(SmartContract)
	school, err := contract.CreateSchool(ctx, "school", "city")
	if err != nil {
		t.Fatal(err)
	}
	stub.SetTransaction("student-1", 101)
	student, err := contract.CreateStudent(ctx, "student")
	if err != nil {
		t.Fatal(err)
	}
	stub.SetTransaction("enroll-1", 102)
	record0, err := contract.EnrollStudent(ctx, school.Address, school.PriKey, student.Address)
	if err != nil || record0.ID != 0 {
		t.Fatalf("enroll record = %+v, %v", record0, err)
	}
	if _, err := contract.EnrollStudent(ctx, school.Address, school.PriKey, student.Address); err == nil {
		t.Fatal("duplicate enrollment must be rejected")
	}
	stub.SetTransaction("graduate-1", 103)
	record1, err := contract.UpdateDiploma(ctx, school.Address, school.PriKey, student.Address, "graduated")
	if err != nil || record1.ID != 1 {
		t.Fatalf("graduation record = %+v, %v", record1, err)
	}
	if _, err := contract.UpdateDiploma(ctx, school.Address, school.PriKey, student.Address, "graduated"); err == nil {
		t.Fatal("duplicate diploma status must be rejected")
	}
	stub.SetTransaction("student-2", 104)
	unenrolled, _ := contract.CreateStudent(ctx, "not-enrolled")
	if _, err := contract.UpdateDiploma(ctx, school.Address, school.PriKey, unenrolled.Address, "graduated"); err == nil {
		t.Fatal("diploma update before enrollment must be rejected")
	}
}
