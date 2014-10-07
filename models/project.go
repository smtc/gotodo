package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/gocache"
	"github.com/smtc/goutils"
)

const (
	PROJECT_CACHE_KEY = "gotodo_project"
)

type Project struct {
	Id         int64  `json:"id"`
	Name       string `sql:"size:100" json:"name"`
	Chief      int64  `json:"chief"`
	Members    string `sql:"size:10000" json:"members"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedBy  int64  `json:"updated_by"`
	UpdatedAt  int64  `json:"updated_at"`
	Level      int    `json:"level"`
	Ongoing    int    `json:"ongoing"`
	Finished   int    `json:"finished"`
	Expired    int    `json:"expired"`
	Visibility int    `json:"visibility"` //0-全部可见 1-仅参与者可见
	Des        string `sql:"size:512" json:"des"`

	ChiefName string  `sql:"-" json:"chief_name"`
	UsersName string  `sql:"-" json:"users_name"`
	LevelText string  `sql:"-" json:"level_text"`
	Users     []int64 `sql:"-" json:"users"`
}

func getProjectDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (p *Project) Refresh() error {
	db := getProjectDB()
	return db.First(p, p.Id).Error
}

func GetProjectList() ([]Project, error) {
	var projects = []Project{}

	m, err := GetProjectCache()
	if err != nil {
		return projects, err
	}

	for _, p := range m {
		projects = append(projects, p)
	}

	return projects, nil
}

func getProjectList() ([]Project, error) {
	var (
		db       = getProjectDB()
		projects []Project
		err      error
	)

	/*
		err = db.Model(&Project{}).Where(where).Count(&total).Error
		if err != nil {
			return 0, nil, err
		}

		if total == 0 {
			return 0, []Project{}, nil
		}
		err = db.Where(where).Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&projects).Error
	*/

	err = db.Find(&projects).Error
	var p *Project
	for i := 0; i < len(projects); i++ {
		p = &projects[i]
		p.LevelText = LEVELS[p.Level]
		p.setUsers()
	}

	return projects, err
}

func (p *Project) setUsers() {
	p.UsersName = GetUserName(p.Chief) + "(*), "
	p.UsersName += GetMultUserName(p.Members)

	users := []int64{p.Chief}
	for _, s := range strings.Split(p.Members, ",") {
		users = append(users, goutils.ToInt64(s, 0))
	}
	p.Users = users
}

func (p *Project) HasMember(id int64) bool {
	for _, u := range p.Users {
		if u == id {
			return true
		}
	}
	return false
}

func (p *Project) Save() error {
	var (
		db    = getProjectDB()
		ids   []string
		chief string
		err   error
	)

	if p.Id == 0 {
		p.CreatedAt = time.Now().Unix()
	}

	for _, id := range strings.Split(p.Members, ",") {
		chief = fmt.Sprintf("%v", p.Chief)
		if chief != id {
			ids = append(ids, id)
		}
	}
	p.Members = strings.Join(ids, ",")

	p.setUsers()
	p.UpdatedAt = time.Now().Unix()

	err = db.Save(p).Error
	setProjectCache(p)

	return err
}

func ProjectDelete(id int64) error {
	var (
		db    = getProjectDB()
		count int
		err   error
	)

	err = db.Model(&Task{}).Where("project_id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Project has %d tasks, can't delete.", count)
	}

	return db.Where("id = ?", id).Delete(&Project{}).Error
}

func GetProject(id int64) (*Project, error) {
	projects, err := GetProjectCache()
	if err != nil {
		return nil, err
	}

	if p, ok := projects[id]; ok {
		return &p, nil
	} else {
		p.Id = id
		err = p.Refresh()
		return &p, err
	}

}

func GetProjectName(id int64) string {
	p, err := GetProject(id)
	if err != nil {
		return ""
	}

	return p.Name
}

func setProjectCache(p *Project) {
	projects, _ := GetProjectCache()
	projects[p.Id] = *p

	cache := gocache.GetCache()
	cache.Set(PROJECT_CACHE_KEY, projects, 0)
}

func GetProjectCache() (map[int64]Project, error) {
	var (
		cache    = gocache.GetCache()
		pc       = map[int64]Project{}
		projects []Project
		err      error
	)

	v, suc := cache.Get(PROJECT_CACHE_KEY)
	if suc {
		pc = v.(map[int64]Project)
	} else {
		projects, err = getProjectList()
		if err != nil {
			return nil, err
		}
		for _, p := range projects {
			pc[p.Id] = p
		}
		cache.Set(PROJECT_CACHE_KEY, pc, 0)
	}
	return pc, nil
}
