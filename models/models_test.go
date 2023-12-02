package models_test

import (
	"testing"

	"example.com/webservice/models"

	"errors"
)

type MockDB struct {
	TestPersons []models.Person
}

func (m *MockDB) GetPersons(count int) ([]models.Person, error) {
	return m.TestPersons, nil
}

func (m *MockDB) AddPerson(person models.Person) error {
	m.TestPersons = append(m.TestPersons, person)
	return nil
}

func (m *MockDB) UpdatePerson(person models.Person) error {
	for i, p := range m.TestPersons {
		if p.Id == person.Id {
			m.TestPersons[i] = person
			return nil
		}
	}
	return errors.New("Kayıt bulunamadı")
}

func (m *MockDB) DeletePerson(id int) error {
	for i, p := range m.TestPersons {
		if p.Id == id {
			m.TestPersons = append(m.TestPersons[:i], m.TestPersons[i+1:]...)
			return nil
		}
	}
	return errors.New("Kayıt bulunamadı")
}

func TestCRUDOperations(t *testing.T) {
	// Test veritabanı olarak mock oluşturuldu
	mockDB := &MockDB{}

	// AddPerson testi
	newPerson := models.Person{Id: 1, FirstName: "Ali", LastName: "Veli", Email: "aliveli@test.com", IpAddress: "192.168.1.1"}
	err := mockDB.AddPerson(newPerson)
	if err != nil {
		t.Errorf("Kişi eklenirken hata oluştu: %v", err)
	}

	t.Logf("Kişi eklendi: %+v", newPerson)

	// GetPersons testi
	persons, err := mockDB.GetPersons(10)
	if err != nil {
		t.Errorf("Kişiler alınırken hata oluştu: %v", err)
	}

	expectedCount := 1 // Beklenen kişi sayısı
	if len(persons) != expectedCount {
		t.Errorf("Beklenen kişi sayısı alınmadı. Beklenen: %d, Alınan: %d", expectedCount, len(persons))
	}

	t.Logf("Alınan kişiler: %+v", persons)

	// UpdatePerson testi
	updatePerson := models.Person{Id: 1, FirstName: "Harry", LastName: "Potter", Email: "harrypotter@test2.com", IpAddress: "192.168.1.2"}
	err = mockDB.UpdatePerson(updatePerson)
	if err != nil {
		t.Errorf("Kişi güncellenirken hata oluştu: %v", err)
	}

	t.Logf("Kişi güncellendi: %+v", updatePerson)

	// GetPersons Güncelleme sonrası ikinci test
	persons, err = mockDB.GetPersons(10)
	if err != nil {
		t.Errorf("Kişiler alınırken hata oluştu: %v", err)
	}

	expectedCount = 1 // Beklenen kişi sayısı
	if len(persons) != expectedCount {
		t.Errorf("Beklenen kişi sayısı alınmadı. Beklenen: %d, Alınan: %d", expectedCount, len(persons))
	}

	t.Logf("Güncellemeden Sonra Alınan kişiler: %+v", persons)

	// DeletePerson testi
	err = mockDB.DeletePerson(1)
	if err != nil {
		t.Errorf("Kişi silinirken hata oluştu: %v", err)
	}

	t.Logf("Kişi silindi: ID=%d", 1)

	// Tekrar GetPersons testi
	persons, err = mockDB.GetPersons(10)
	if err != nil {
		t.Errorf("Kişiler alınırken hata oluştu: %v", err)
	}

	expectedCount = 0 // Beklenen kişi sayısı 0 (silindiği için)
	if len(persons) != expectedCount {
		t.Errorf("Beklenen kişi sayısı alınmadı. Beklenen: %d, Alınan: %d", expectedCount, len(persons))
	}

	t.Logf("Tekrarlanan getPersons testinde alınan kişiler: %+v", persons)
}

type MockUserDB struct {
	TestUsers []models.User
}

func (m *MockUserDB) GetUsers() ([]models.User, error) {
	return m.TestUsers, nil
}

func (m *MockUserDB) AddUser(user models.User) error {
	m.TestUsers = append(m.TestUsers, user)
	return nil
}

func (m *MockUserDB) UpdateUser(user models.User) error {
	for i, u := range m.TestUsers {
		if u.ID == user.ID {
			m.TestUsers[i] = user
			return nil
		}
	}
	return errors.New("Kullanıcı bulunamadı")
}

func (m *MockUserDB) DeleteUser(id int) error {
	for i, u := range m.TestUsers {
		if u.ID == id {
			m.TestUsers = append(m.TestUsers[:i], m.TestUsers[i+1:]...)
			return nil
		}
	}
	return errors.New("Kullanıcı bulunamadı")
}

func TestCRUDOperationsForUser(t *testing.T) {
	mockDB := &MockUserDB{}

	// AddUser testi
	newUser := models.User{ID: 1, Username: "testuser", Email: "testuser@example.com", Password: "test1234"}
	err := mockDB.AddUser(newUser)
	if err != nil {
		t.Errorf("Kullanıcı eklenirken hata oluştu: %v", err)
	}

	t.Logf("Kullanıcı eklendi: %+v", newUser)

	// GetUsers testi
	users, err := mockDB.GetUsers()
	if err != nil {
		t.Errorf("Kullanıcılar alınırken hata oluştu: %v", err)
	}

	expectedCount := 1 // Beklenen kullanıcı sayısı
	if len(users) != expectedCount {
		t.Errorf("Beklenen kullanıcı sayısı alınmadı. Beklenen: %d, Alınan: %d", expectedCount, len(users))
	}

	t.Logf("Alınan kullanıcılar: %+v", users)

	// UpdateUser testi
	updateUser := models.User{ID: 1, Username: "testuserGÜNCELLENDİ", Email: "testuserGÜNCELLENDİ@example.com", Password: "test1234GÜNCELLENDİ"}
	err = mockDB.UpdateUser(updateUser)
	if err != nil {
		t.Errorf("Kullanıcı güncellenirken hata oluştu: %v", err)
	}

	t.Logf("Kullanıcı güncellendi: %+v", updateUser)

	// GetUsers Güncelleme sonrası ikinci test
	users, err = mockDB.GetUsers()
	if err != nil {
		t.Errorf("Kullanıcılar alınırken hata oluştu: %v", err)
	}

	expectedCount = 1 // Beklenen kullanıcı sayısı
	if len(users) != expectedCount {
		t.Errorf("Beklenen kullanıcı sayısı alınmadı. Beklenen: %d, Alınan: %d", expectedCount, len(users))
	}

	t.Logf("Güncellemeden Sonra Alınan kullanıcılar: %+v", users)

	// DeleteUser testi
	err = mockDB.DeleteUser(1)
	if err != nil {
		t.Errorf("Kullanıcı silinirken hata oluştu: %v", err)
	}

	t.Logf("Kullanıcı silindi: ID=%d", 1)

	// Tekrar GetUsers testi
	users, err = mockDB.GetUsers()
	if err != nil {
		t.Errorf("Kullanıcılar alınırken hata oluştu: %v", err)
	}

	expectedCount = 0 // Beklenen kullanıcı sayısı 0 (silindiği için)
	if len(users) != expectedCount {
		t.Errorf("Beklenen kullanıcı sayısı alınmadı. Beklenen: %d, Alınan: %d", expectedCount, len(users))
	}

	t.Logf("Tekrarlanan GetUsers testinde alınan kullanıcılar: %+v", users)
}
