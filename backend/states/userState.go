package states

import (
	"go-react/backend/models"
	"sync"
)

var (
	CurrentUser models.User
	UserMutex   sync.RWMutex
)