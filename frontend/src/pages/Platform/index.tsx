import { useEffect, useMemo, useState } from "react";
import { Box, Flex, Heading, IconButton, Text } from "@radix-ui/themes";
import { Edit, Plus } from "lucide-react";

import CButton from "../../components/common/CButton";
import CInput from "../../components/common/CInput";
import CModel from "../../components/common/CModel";
import CTable, { type Column } from "../../components/common/CTable";
import { useErrorHandler } from "../../hooks/useErrorHandler";
import PlatformService, {
  type PlatformTenant,
  type PlatformTenantInput,
} from "../../services/platform/platform.service";

type TenantFormState = PlatformTenantInput & {
  EditID?: number;
};

const emptyForm: TenantFormState = {
  Tenant: {
    Name: "",
    Code: "",
    SubDomain: "",
  },
  Location: {
    Code: "LOC_1",
    Address: "",
    City: "",
    State: "",
    Country: "India",
    ZipCode: "",
  },
  Admin: {
    Email: "",
    Password: "",
    Name: "",
  },
};

function isCompleteCreateForm(form: TenantFormState): boolean {
  return Boolean(
    form.Tenant.Name.trim() &&
    form.Tenant.Code.trim() &&
    form.Tenant.SubDomain.trim() &&
    form.Location.Code.trim() &&
    form.Location.Address.trim() &&
    form.Location.City.trim() &&
    form.Location.State.trim() &&
    form.Location.Country.trim() &&
    form.Location.ZipCode.trim() &&
    form.Admin.Email.trim() &&
    form.Admin.Password.trim() &&
    form.Admin.Name.trim()
  );
}

function PlatformLogin({ onLogin }: { onLogin: () => void }) {
  const { showError } = useErrorHandler();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async () => {
    try {
      setLoading(true);
      await PlatformService.login(username, password);
      onLogin();
    } catch (err) {
      showError(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Flex
      align="center"
      justify="center"
      style={{ minHeight: "100vh", background: "var(--gray-2)" }}
    >
      <Box
        style={{
          width: 360,
          padding: 24,
          border: "1px solid var(--gray-6)",
          borderRadius: 6,
          background: "white",
        }}
      >
        <Flex direction="column" gap="4">
          <Box>
            <Heading size="5">Platform Admin</Heading>
            <Text size="2" color="gray">
              Sign in to provision and manage tenants.
            </Text>
          </Box>
          <CInput
            label="Username"
            value={username}
            onChange={setUsername}
            fullWidth
          />
          <CInput
            label="Password"
            type="password"
            value={password}
            onChange={setPassword}
            fullWidth
          />
          <CButton
            label="Login"
            fullWidth
            processing={loading}
            disabled={!username || !password}
            onClick={submit}
          />
        </Flex>
      </Box>
    </Flex>
  );
}

function TenantDrawer({
  form,
  onChange,
  onClose,
  onSubmit,
  saving,
}: {
  form: TenantFormState;
  onChange: (form: TenantFormState) => void;
  onClose: () => void;
  onSubmit: () => void;
  saving: boolean;
}) {
  const isEdit = Boolean(form.EditID);

  const setTenant = (key: keyof TenantFormState["Tenant"], value: string) =>
    onChange({ ...form, Tenant: { ...form.Tenant, [key]: value } });
  const setLocation = (key: keyof TenantFormState["Location"], value: string) =>
    onChange({ ...form, Location: { ...form.Location, [key]: value } });
  const setAdmin = (key: keyof TenantFormState["Admin"], value: string) =>
    onChange({ ...form, Admin: { ...form.Admin, [key]: value } });

  return (
    <CModel
      open
      size="md"
      title={isEdit ? "Edit Tenant" : "Provision Tenant"}
      onClose={onClose}
      actions={
        <Flex gap="3">
          <CButton label="Cancel" variant="soft" onClick={onClose} />
          <CButton
            label={isEdit ? "Update" : "Create"}
            onClick={onSubmit}
            processing={saving}
            disabled={saving || (!isEdit && !isCompleteCreateForm(form))}
          />
        </Flex>
      }
    >
      <Flex direction="column" gap="5">
        <Flex direction="column" gap="3">
          <Text size="2" weight="bold">
            Tenant
          </Text>
          <CInput
            label="Name"
            value={form.Tenant.Name}
            onChange={(value) => setTenant("Name", value)}
            fullWidth
          />
          <CInput
            label="Code"
            value={form.Tenant.Code}
            onChange={(value) => setTenant("Code", value)}
            disabled={isEdit}
            fullWidth
          />
          <CInput
            label="Subdomain"
            value={form.Tenant.SubDomain}
            onChange={(value) => setTenant("SubDomain", value)}
            disabled={isEdit}
            fullWidth
          />
        </Flex>

        {!isEdit && (
          <>
            <Flex direction="column" gap="3">
              <Text size="2" weight="bold">
                Primary Location
              </Text>
              <CInput
                label="Code"
                value={form.Location.Code}
                onChange={(value) => setLocation("Code", value)}
                fullWidth
              />
              <CInput
                label="Address"
                value={form.Location.Address}
                onChange={(value) => setLocation("Address", value)}
                fullWidth
              />
              <Flex gap="3">
                <CInput
                  label="City"
                  value={form.Location.City}
                  onChange={(value) => setLocation("City", value)}
                  fullWidth
                />
                <CInput
                  label="State"
                  value={form.Location.State}
                  onChange={(value) => setLocation("State", value)}
                  fullWidth
                />
              </Flex>
              <Flex gap="3">
                <CInput
                  label="Country"
                  value={form.Location.Country}
                  onChange={(value) => setLocation("Country", value)}
                  fullWidth
                />
                <CInput
                  label="Zip Code"
                  value={form.Location.ZipCode}
                  onChange={(value) => setLocation("ZipCode", value)}
                  fullWidth
                />
              </Flex>
            </Flex>

            <Flex direction="column" gap="3">
              <Text size="2" weight="bold">
                Tenant Admin
              </Text>
              <CInput
                label="Name"
                value={form.Admin.Name}
                onChange={(value) => setAdmin("Name", value)}
                fullWidth
              />
              <CInput
                label="Email"
                type="email"
                value={form.Admin.Email}
                onChange={(value) => setAdmin("Email", value)}
                fullWidth
              />
              <CInput
                label="Password"
                type="password"
                value={form.Admin.Password}
                onChange={(value) => setAdmin("Password", value)}
                fullWidth
              />
            </Flex>
          </>
        )}
      </Flex>
    </CModel>
  );
}

function PlatformDashboard({ onLogout }: { onLogout: () => void }) {
  const { showError, showSuccess } = useErrorHandler();
  const [tenants, setTenants] = useState<PlatformTenant[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState<TenantFormState | null>(null);

  const loadTenants = async () => {
    try {
      setLoading(true);
      setTenants(await PlatformService.listTenants());
    } catch (err) {
      showError(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTenants();
  }, []);

  const filteredTenants = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return tenants;
    return tenants.filter((tenant) =>
      [tenant.Name, tenant.Code, tenant.SubDomain]
        .join(" ")
        .toLowerCase()
        .includes(q)
    );
  }, [search, tenants]);

  const submitForm = async () => {
    if (!form) return;
    try {
      setSaving(true);
      if (form.EditID) {
        await PlatformService.updateTenant(form.EditID, {
          Name: form.Tenant.Name,
        });
        showSuccess("Tenant updated successfully");
      } else {
        await PlatformService.createTenant(form);
        showSuccess("Tenant provisioned successfully");
      }
      setForm(null);
      await loadTenants();
    } catch (err) {
      showError(err);
    } finally {
      setSaving(false);
    }
  };

  const columns: Column<PlatformTenant>[] = [
    { key: "Name", label: "Name" },
    { key: "Code", label: "Code" },
    {
      key: "SubDomain",
      label: "Subdomain",
      render: (row) => <Text>{row.SubDomain}</Text>,
    },
    {
      key: "Actions",
      label: "Actions",
      render: (row) => (
        <IconButton
          variant="ghost"
          radius="full"
          style={{ cursor: "pointer" }}
          onClick={() =>
            setForm({
              ...emptyForm,
              EditID: row.ID,
              Tenant: {
                Name: row.Name,
                Code: row.Code,
                SubDomain: row.SubDomain,
              },
            })
          }
        >
          <Edit size={16} />
        </IconButton>
      ),
    },
  ];

  return (
    <Box style={{ minHeight: "100vh", background: "var(--gray-1)" }}>
      <Flex
        justify="between"
        align="center"
        p="4"
        style={{ borderBottom: "1px solid var(--gray-5)", background: "white" }}
      >
        <Box>
          <Heading size="5">Platform Tenants</Heading>
          <Text size="2" color="gray">
            Provision tenants with primary location and admin user.
          </Text>
        </Box>
        <CButton label="Logout" variant="soft" onClick={onLogout} />
      </Flex>

      <Box>
        <CTable
          title="Tenants"
          data={filteredTenants}
          rowKey="ID"
          columns={columns}
          loading={loading}
          searchable
          searchPlaceholder="Search tenants"
          onSearch={setSearch}
          onRefresh={loadTenants}
          actions={
            <CButton
              label="Create Tenant"
              startIcon={<Plus size={16} />}
              onClick={() =>
                setForm(JSON.parse(JSON.stringify(emptyForm)) as TenantFormState)
              }
            />
          }
        />
      </Box>

      {form && (
        <TenantDrawer
          form={form}
          onChange={setForm}
          onClose={() => setForm(null)}
          onSubmit={submitForm}
          saving={saving}
        />
      )}
    </Box>
  );
}

export default function PlatformPage() {
  const [checking, setChecking] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);

  const checkSession = async () => {
    try {
      await PlatformService.me();
      setAuthenticated(true);
    } catch {
      setAuthenticated(false);
    } finally {
      setChecking(false);
    }
  };

  useEffect(() => {
    checkSession();
  }, []);

  const logout = async () => {
    await PlatformService.logout();
    setAuthenticated(false);
  };

  if (checking) return null;
  if (!authenticated) return <PlatformLogin onLogin={checkSession} />;
  return <PlatformDashboard onLogout={logout} />;
}
