package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/router/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) RBACRolesPermsGrid(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	grid, err := h.SS.RBAC.GetAllRolePermissions(tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Role permission grid fetched successfully.",
		"Data": gin.H{
			"Roles":       h.SS.RBAC.GetAllRoles(tenantID),
			"Permissions": h.SS.RBAC.GetAllPermisssions(),
			"Grid":        grid,
		},
	})

}
func (h *Handlers) RBACAddRole(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	addRole := &dto.AddRole{
		TenantID: tenantID,
	}
	if err := c.ShouldBindBodyWithJSON(addRole); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	err := h.SS.RBAC.AddUserRole(addRole)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Role created successfully.",
	})
}

func (h *Handlers) RBACUpdateRolePerms(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	updateRoles := &dto.UpdateRoles{
		TenantID: tenantID,
	}
	if err := c.ShouldBindBodyWithJSON(updateRoles); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	err := h.SS.RBAC.UpdateUserRole(updateRoles)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Role permissions updated successfully.",
	})
}
