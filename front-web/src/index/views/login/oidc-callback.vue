<template>
    <div>
        <span>&nbsp;</span>
    </div>
</template>

<script>
import { mapActions } from "vuex"

export default {
    name: "OidcCallback",
    created() {
        this.oidcSignInCallback()
            .then(redirectPath => {
                console.log(redirectPath)
                // this.$router.push({ path: redirectPath })
                // this.$router.push({ path: redirectPath || '/' })
                this.$router.push(redirectPath)
            })
            .catch(err => {
                console.error(err)
                this.$router.push("/oidc-callback-error") // Handle errors any way you want
            })
    },
    methods: {
        ...mapActions("oidc", ["oidcSignInCallback"])
    }
}
</script>
