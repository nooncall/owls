<template>
  <div id="userLayout">
    <div class="login_panle">
      <div class="login_panle_form">
        <div class="login_panle_form_title">
          <img
            class="login_panle_form_title_logo"
            :src="$GIN_VUE_ADMIN.appLogo"
            alt
          >
          <p class="login_panle_form_title_p">{{ $GIN_VUE_ADMIN.appName }}</p>
        </div>
        <el-form
          ref="loginForm"
          :model="loginFormData"
          :rules="rules"
          @keyup.enter="submitForm"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginFormData.username"
              placeholder="请输入用户名"
            >
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <user />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginFormData.password"
              :type="lock === 'lock' ? 'password' : 'text'"
              placeholder="请输入密码"
            >
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <component
                      :is="lock"
                      @click="changeLock"
                    />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="captcha">
            <div class="vPicBox">
              <el-input
                v-model="loginFormData.captcha"
                placeholder="请输入验证码"
                style="width: 60%"
              />
              <div class="vPic">
                <img
                  v-if="picPath"
                  :src="picPath"
                  alt="请输入验证码"
                  @click="loginVerify()"
                >
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              style="width: 22%"
              size="large"
              @click="checkInit"
            >前往初始化</el-button>
            <el-button
              v-if="loginType == 'registry'"
              type="primary"
              size="large"
              style="width: 22%; margin-left: 28%"
              @click="addUser"
            >注 册</el-button>
            <el-button
                type="primary"
                size="large"
                style="width: 22%; float: right; margin-left: 6%"
                @click="submitForm"
            >登 录</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="login_panle_right" />
      <div class="login_panle_foot">
        <div class="links">
          <a href="http://doc.henrongyi.top/" target="_blank">
            <img src="@/assets/docs.png" class="link-icon">
          </a>
          <a href="https://support.qq.com/product/371961" target="_blank">
            <img src="@/assets/kefu.png" class="link-icon">
          </a>
          <a
            href="https://github.com/flipped-aurora/gin-vue-admin"
            target="_blank"
          >
            <img src="@/assets/github.png" class="link-icon">
          </a>
          <a href="https://space.bilibili.com/322210472" target="_blank">
            <img src="@/assets/video.png" class="link-icon">
          </a>
        </div>
        <div class="copyright">
          <bootomInfo />
        </div>
      </div>
    </div>

    <el-dialog
        v-model="addUserDialog"
        custom-class="user-dialog"
        title="注册"
        :show-close="false"
        :close-on-press-escape="false"
        :close-on-click-modal="false"
    >
      <div style="height:45vh;overflow:auto;padding:0 12px;">
        <el-form ref="userForm" :rules="registryRules" :model="userInfo" label-width="80px">
          <el-form-item label="用户名" prop="userName">
            <el-input v-model="userInfo.userName" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="userInfo.password" />
          </el-form-item>
          <el-form-item label="昵称" prop="nickName">
            <el-input v-model="userInfo.nickName" />
          </el-form-item>
          <el-form-item label="手机号" prop="phone">
            <el-input v-model="userInfo.phone" />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="userInfo.email" />
          </el-form-item>
        </el-form>

      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeAddUserDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterAddUserDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script>
export default {
  name: 'Login',
}
</script>

<script setup>
import { captcha, baseRegister, loginModel } from '@/api/user'
import { checkDB } from '@/api/initdb'
import bootomInfo from '@/view/layout/bottomInfo/bottomInfo.vue'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
const router = useRouter()
// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 3) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}

const loginType = ref('')
const getLoginType = async() => {
  const model = await loginModel()
  loginType.value = model.data
}
getLoginType()

// 获取验证码
const loginVerify = () => {
  captcha({}).then((ele) => {
    rules.captcha[1].max = ele.data.captchaLength
    rules.captcha[1].min = ele.data.captchaLength
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
  })
}
loginVerify()

// 登录相关操作
const lock = ref('lock')
const changeLock = () => {
  lock.value = lock.value === 'lock' ? 'unlock' : 'lock'
}

const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: 'admin',
  password: '123456',
  captcha: '',
  captchaId: '',
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    {
      message: '验证码格式不正确',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const login = async() => {
  return await userStore.LoginIn(loginFormData)
}
const submitForm = () => {
  loginForm.value.validate(async(v) => {
    if (v) {
      const flag = await login()
      if (!flag) {
        loginVerify()
      }
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      loginVerify()
      return false
    }
  })
}

// 注册

// 弹窗相关
const userInfo = ref({
  username: '',
  password: '',
  nickName: '',
  headerImg: '',
  authorityId: '',
  authorityIds: [],
})

const registryRules = ref({
  userName: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 5, message: '最低5位字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入用户密码', trigger: 'blur' },
    { min: 6, message: '最低6位字符', trigger: 'blur' }
  ],
  nickName: [
    { required: true, message: '请输入用户昵称', trigger: 'blur' }
  ],
})

const userForm = ref(null)
const enterAddUserDialog = async() => {
  userInfo.value.authorityId = userInfo.value.authorityIds[0]
  userForm.value.validate(async valid => {
    if (valid) {
      const req = {
        ...userInfo.value
      }
      const res = await baseRegister(req)
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '创建成功' })
        closeAddUserDialog()
      }
    }
  })
}

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
  userForm.value.resetFields()
  userInfo.value.headerImg = ''
  addUserDialog.value = false
}

const addUser = () => {
  addUserDialog.value = true
}

// 跳转初始化
const checkInit = async() => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      router.push({ name: 'Init' })
    } else {
      ElMessage({
        type: 'info',
        message: '已配置数据库信息，无法初始化',
      })
    }
  }
}

</script>

<style lang="scss" scoped>
@import "@/style/newLogin.scss";
</style>
