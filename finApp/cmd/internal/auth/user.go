package auth

func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []model.User
		db.Order("id").Find(&us)
		c.JSON(200, us)
	}
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in struct {Name, Email, Password, Role string}
		if c.ShouldBindJSON(&in) != nil || len(in.Password) < 8 || (in.Role != "admin" && in.Role != "user") {
			c.JSON(400, gin.H{"error": "nama, email, password (min 8 karakter), dan role (admin/user) harus diisi"})
			return
		}
		h, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		u := model.User{Name: in.Name, Email: in.Email, PasswordHash: string(h), Role: in.Role, IsActive: true}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(409, gin.H{"error": "Email sudah digunakan"})
			return
		}
		c.JSON(201, u)
	}
}

func PatchUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in struct {
		Role  *string `json:"role"`
		IsActive *bool   `json:"is_active"`	
		}
	if c.ShouldBindJSON(&in) != nil {
		c.JSON(400, gin.H{"error": "Input Tidak Valid"})
		return
	}
	upd := map[string]any{}
	if in.Role != nil && (*in.Role == "admin" || *in.Role == "user") {
		upd["role"] = *in.Role
	}
	if in.IsActive != nil {
		upd["is_active"] = *in.IsActive
	}
	db.Model(&model.User{}).Where("id = ?", c.Param("id")).Updates(upd)
	c.Status(204)
	}
}