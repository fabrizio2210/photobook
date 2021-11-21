<template>
    <div>
        <h3 v-if="project" >{{ project.name }}</h3>
        <div id="nav">
            <router-link v-if="project" :to="'/projects/' + project.id">Executions</router-link> |
            <router-link v-if="project" :to="'/projects/' + project.id + '/environments'">Environment</router-link> |
            <router-link v-if="project" :to="'/projects/' + project.id + '/settings'">Settings</router-link>
        </div>
    <router-view />
    </div>
</template>

<script>

export default {
    data () {
        return {
        }
    },
    computed: {
        project () {
            const project_id = this.$route.params.project_id;
            return this.$store.state.projects.all.projects_dict[project_id];
        }
    },
    created () {
        const project_id = this.$route.params.project_id;
        if (!this.project){
            this.$store.dispatch('projects/get', { project_id });
        }
    }
};
</script>
