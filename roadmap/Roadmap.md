| #   | Feature / Route                                                  | Allowed Roles                | Status     | Notes / To Do                             |
| --- | ---------------------------------------------------------------- | ---------------------------- | ---------- | ----------------------------------------- |
| 1   | Health Check (`/health`)                                         | Public                       | ✅ Done    | Basic server health endpoint              |
| 2   | Auth - Me (`/auth/v1/me`)                                        | Authenticated Users          | ✅ Done    | Returns logged-in user info               |
| 3   | Auth - Login (`/auth/v1/login`)                                  | Public                       | ✅ Done    | Email/password login                      |
| 4   | Auth - Logout (`/auth/v1/logout`)                                | Authenticated Users          | ✅ Done    | Clears session / token                    |
| 5   | OAuth - Google Login (`/auth/v1/oauth/login/google`)             | Public                       | ✅ Done    | Redirect to Google OAuth                  |
| 6   | OAuth - Google Callback (`/auth/v1/oauth/login/google/callback`) | Public                       | ✅ Done    | Handles Google OAuth callback             |
| 7   | OAuth - GitHub Login (`/auth/v1/oauth/login/github`)             | Public                       | ✅ Done    | Redirect to GitHub OAuth                  |
| 8   | OAuth - GitHub Callback (`/auth/v1/oauth/login/github/callback`) | Public                       | ✅ Done    | Handles GitHub OAuth callback             |
| 9   | File Upload (`/api/v1/files/upload`)                             | Tenant Admin, Employee       | ✅ Done    | Upload file endpoint                      |
| 10  | Get File by ID (`/api/v1/files/get/:id`)                         | Tenant Admin, Employee       | ✅ Done    | Download / retrieve file                  |
| 11  | Tenant - Create (`/api/v1/tenant/create`)                        | Super Admin                  | ✅ Done    | Create a new tenant                       |
| 12  | Tenant - Update (`/api/v1/tenant/update`)                        | Super Admin, Tenant Admin    | ✅ Done    | Update tenant info                        |
| 13  | Tenant - Get by ID (`/api/v1/tenant/:tenant_id`)                 | Super Admin, Tenant Admin    | ✅ Done    | Retrieve tenant info                      |
| 14  | Tenant - Get All (`/api/v1/tenant/all`)                          | Super Admin                  | ✅ Done    | List all tenants                          |
| 15  | Tenant Location - Create (`/api/v1/location/create`)             | Tenant Admin                 | ✅ Done    | Add tenant location                       |
| 16  | Tenant Location - Update (`/api/v1/location/update`)             | Tenant Admin                 | ✅ Done    | Update tenant location                    |
| 17  | Tenant Location - Get by ID (`/api/v1/location/:LocationID`)     | Tenant Admin, Employee       | ✅ Done    | Retrieve location info                    |
| 18  | Tenant Location - Get All (`/api/v1/location/all`)               | Tenant Admin                 | ✅ Done    | List all tenant locations                 |
| 19  | Employee Signup (`/api/v1/employee/signup`)                      | Tenant Admin                 | ✅ Done    | Create employee under tenant              |
| 20  | Employee Update (`/api/v1/employee/update`)                      | Tenant Admin                 | ✅ Done    | Update employee info                      |
| 21  | Employee Get by ID (`/api/v1/employee/:UserID`)                  | Tenant Admin, Employee       | ✅ Done    | Get employee details                      |
| 22  | Employee Get All (`/api/v1/employee/all`)                        | Tenant Admin                 | ✅ Done    | List all employees                        |
| 23  | User Signup (`/api/v1/user/signup`)                              | Public                       | ✅ Done    | Create individual user                    |
| 24  | User Update (`/api/v1/user/update`)                              | User                         | ✅ Done    | Update user info                          |
| 25  | User Get by ID (`/api/v1/user/:UserID`)                          | User, Tenant Admin           | ✅ Done    | Get user info                             |
| 26  | Menu Item - Create (`/api/v1/item/create`)                       | Tenant Admin                 | ✅ Done    | Add a menu item                           |
| 27  | Menu Item - Update (`/api/v1/item/update`)                       | Tenant Admin                 | ✅ Done    | Update menu item                          |
| 28  | Menu Item - Get by ID (`/api/v1/item/:ItemID`)                   | Employee, Tenant Admin       | ✅ Done    | Retrieve item info                        |
| 29  | Menu Item - Get All (`/api/v1/item/all`)                         | Employee, Tenant Admin       | ✅ Done    | List all menu items                       |
| 30  | Order - Create (`/api/v1/order/create`)                          | Employee, User               | ⬜ Pending | Implement order creation API              |
| 31  | Order - Update Status (`/api/v1/order/update/:id`)               | Tenant Admin, Employee       | ⬜ Pending | Implement order status update             |
| 32  | Order - Get by ID (`/api/v1/order/:id`)                          | User, Employee, Tenant Admin | ⬜ Pending | Retrieve order details                    |
| 33  | Order - Get All (`/api/v1/order/all`)                            | Tenant Admin, Employee       | ⬜ Pending | List all orders                           |
| 34  | Middleware for Protected Routes                                  | Authenticated Routes         | ⬜ Pending | Implement middleware to check JWT/session |
| 35  | Integrate RABC (Role-Based Access Control)                       | All Authenticated Routes     | ⬜ Pending | Enforce per-role permissions dynamically  |
