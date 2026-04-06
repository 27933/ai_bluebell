<template>
  <div class="users-page">
    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-select v-model="filterRole" placeholder="角色" style="width: 120px" @change="handleFilter">
        <el-option label="全部角色" value="all" />
        <el-option label="admin" value="admin" />
        <el-option label="author" value="author" />
        <el-option label="reader" value="reader" />
      </el-select>
      <el-select v-model="filterStatus" placeholder="状态" style="width: 120px" @change="handleFilter">
        <el-option label="全部状态" value="all" />
        <el-option label="正常" value="active" />
        <el-option label="已封禁" value="inactive" />
      </el-select>

      <div v-if="selectedIds.length > 0" class="batch-bar">
        <span>已选 {{ selectedIds.length }} 人</span>
        <el-button size="small" type="danger" @click="batchBan('inactive')">批量封禁</el-button>
        <el-button size="small" type="success" @click="batchBan('active')">批量解封</el-button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <el-table :data="users" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="40" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="100" />
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="roleTagType(row.role)" size="small">{{ row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '正常' : '封禁' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="130">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              text
              :disabled="row.id === currentUserId"
              @click="openRoleDialog(row)"
            >改角色</el-button>
            <el-button
              v-if="row.status === 'active'"
              size="small"
              type="danger"
              text
              :disabled="row.id === currentUserId"
              @click="handleStatusChange(row, 'inactive')"
            >封禁</el-button>
            <el-button
              v-else
              size="small"
              type="success"
              text
              @click="handleStatusChange(row, 'active')"
            >解封</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadUsers"
        />
      </div>
    </div>

    <!-- 修改角色弹窗 -->
    <el-dialog v-model="roleDialogVisible" title="修改角色" width="320px" :close-on-click-modal="false">
      <div style="margin-bottom: 1rem; color: #475569;">
        用户：<strong>{{ editingUser?.username }}</strong>
      </div>
      <el-select v-model="newRole" style="width: 100%">
        <el-option label="admin" value="admin" />
        <el-option label="author" value="author" />
        <el-option label="reader" value="reader" />
      </el-select>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="roleLoading" @click="submitRoleChange">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '../../stores/auth'
import {
  getUserList,
  updateUserRole,
  updateUserStatus,
  batchUpdateUserStatus,
} from '../../api/admin'
import type { AdminUser } from '../../api/admin'

const authStore = useAuthStore()
const currentUserId = computed(() => authStore.user?.id)

const users = ref<AdminUser[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const filterRole = ref('all')
const filterStatus = ref('all')
const selectedIds = ref<string[]>([])

// 修改角色弹窗
const roleDialogVisible = ref(false)
const editingUser = ref<AdminUser | null>(null)
const newRole = ref('')
const roleLoading = ref(false)

onMounted(() => loadUsers())

async function loadUsers() {
  loading.value = true
  try {
    const res = await getUserList({
      role: filterRole.value,
      status: filterStatus.value,
      page: page.value,
      size: pageSize,
    })
    if (res.code === 1000) {
      users.value = res.data.list || []
      total.value = res.data.total
    }
  } catch {
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

function handleFilter() {
  page.value = 1
  loadUsers()
}

function handleSelectionChange(rows: AdminUser[]) {
  selectedIds.value = rows.map(r => r.id)
}

function openRoleDialog(row: AdminUser) {
  editingUser.value = row
  newRole.value = row.role
  roleDialogVisible.value = true
}

async function submitRoleChange() {
  if (!editingUser.value || newRole.value === editingUser.value.role) {
    roleDialogVisible.value = false
    return
  }
  roleLoading.value = true
  try {
    const res = await updateUserRole(editingUser.value.id, newRole.value)
    if (res.code === 1000) {
      ElMessage.success('角色修改成功')
      roleDialogVisible.value = false
      await loadUsers()
    } else {
      ElMessage.error('修改失败')
    }
  } catch {
    ElMessage.error('操作失败')
  } finally {
    roleLoading.value = false
  }
}

async function handleStatusChange(row: AdminUser, status: string) {
  if (row.id === currentUserId.value && status === 'inactive') {
    ElMessage.warning('不能封禁自己')
    return
  }
  const action = status === 'inactive' ? '封禁' : '解封'
  try {
    await ElMessageBox.confirm(`确认${action}用户 ${row.username}？`, `确认${action}`, {
      type: 'warning',
    })
    const res = await updateUserStatus(row.id, status)
    if (res.code === 1000) {
      ElMessage.success(`${action}成功`)
      await loadUsers()
    }
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

async function batchBan(status: string) {
  const action = status === 'inactive' ? '封禁' : '解封'
  try {
    await ElMessageBox.confirm(`确认批量${action}选中的 ${selectedIds.value.length} 个用户？`, `批量${action}`, {
      type: 'warning',
    })
    const res = await batchUpdateUserStatus(selectedIds.value, status)
    if (res.code === 1000) {
      ElMessage.success(`批量${action}成功`)
      loadUsers()
    }
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

function roleTagType(role: string): 'danger' | 'warning' | 'info' {
  if (role === 'admin') return 'danger'
  if (role === 'author') return 'warning'
  return 'info'
}

function formatDate(str: string): string {
  return str ? new Date(str).toLocaleDateString('zh-CN') : '-'
}
</script>

<style scoped>
.users-page {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.batch-bar {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #475569;
  font-size: 0.9rem;
}

.table-card {
  background: white;
  border-radius: 10px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.07);
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}
</style>
