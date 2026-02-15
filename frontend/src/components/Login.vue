<template>
  <v-container fluid class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12 rounded-lg" :dark="darkMode">
          <v-card-title class="text-h4 font-weight-light pa-8 pb-4">
            Welcome Back
          </v-card-title>
          <v-card-subtitle class="text-body-1 pa-8 pt-0 pb-6">
            Sign in to continue
          </v-card-subtitle>
          <v-card-text class="pa-8 pt-0">
            <v-form ref="form" v-model="valid">
              <v-text-field
                v-model="username"
                :rules="usernameRules"
                label="Username"
                prepend-inner-icon="mdi-account"
                outlined
                rounded
                required
                class="mb-4"
              ></v-text-field>

              <v-text-field
                v-model="password"
                :rules="passwordRules"
                :type="showPassword ? 'text' : 'password'"
                label="Password"
                prepend-inner-icon="mdi-lock"
                :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                @click:append="showPassword = !showPassword"
                outlined
                rounded
                required
                class="mb-4"
                @keyup.enter="handleLogin"
              ></v-text-field>

              <v-btn
                :disabled="!valid || loading"
                :loading="loading"
                color="primary"
                large
                rounded
                block
                class="mt-4"
                @click="handleLogin"
              >
                Sign In
              </v-btn>
            </v-form>
          </v-card-text>
          <v-card-actions class="pa-8 pt-0">
            <v-spacer></v-spacer>
            <span class="text-body-2">Don't have an account?</span>
            <v-btn text color="primary" @click="goToSignup">Sign Up</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="4000" top>
      {{ snackbarText }}
      <template v-slot:action="{ attrs }">
        <v-btn text v-bind="attrs" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script>
import EventService from '../services/EventService'

export default {
  name: 'Login',
  data() {
    return {
      valid: false,
      username: '',
      password: '',
      showPassword: false,
      loading: false,
      snackbar: false,
      snackbarText: '',
      snackbarColor: 'error',
      usernameRules: [
        v => !!v || 'Username is required',
        v => (v && v.length >= 3) || 'Username must be at least 3 characters'
      ],
      passwordRules: [
        v => !!v || 'Password is required',
        v => (v && v.length >= 6) || 'Password must be at least 6 characters'
      ]
    }
  },
  computed: {
    darkMode() {
      return this.$vuetify.theme.dark
    }
  },
  methods: {
    handleLogin() {
      if (!this.$refs.form.validate()) {
        return
      }

      this.loading = true

      const lData = {
        username: this.username,
        password: this.password
      }

      EventService.login(lData)
        .then((lRes) => {
          if (lRes.data.status === 's') {
            const lToken = lRes.data.data.token
            const lUser = lRes.data.data.user
            localStorage.setItem('token', lToken)
            localStorage.setItem('user', JSON.stringify(lUser))
            this.showSnackbar('Login successful', 'success')
            setTimeout(() => {
              this.$router.push('/todos')
            }, 500)
          } else {
            this.showSnackbar(lRes.data.message || 'Login failed', 'error')
          }
        })
        .catch((lErr) => {
          if (lErr.response && lErr.response.data && lErr.response.data.message) {
            this.showSnackbar(lErr.response.data.message, 'error')
          } else {
            this.showSnackbar('An error occurred', 'error')
          }
        })
        .finally(() => {
          this.loading = false
        })
    },
    goToSignup() {
      this.$router.push('/signup')
    },
    showSnackbar(pText, pColor) {
      this.snackbarText = pText
      this.snackbarColor = pColor
      this.snackbar = true
    }
  }
}
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
}
</style>

