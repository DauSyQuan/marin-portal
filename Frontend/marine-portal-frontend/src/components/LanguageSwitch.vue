<script setup>
import { computed } from "vue";
import { useI18nStore } from "@/stores/i18n";

const props = defineProps({
  variant: {
    type: String,
    default: "default", // default | compact
  },
});

const i18n = useI18nStore();

const cls = computed(() => ({
  "lang-switch": true,
  "is-compact": props.variant === "compact",
}));
</script>

<template>
  <div :class="cls">
    <button
      class="lang-btn"
      :class="{ active: i18n.lang === 'vi' }"
      @click="i18n.setLang('vi')"
      type="button"
    >
      VN <span class="code">VI</span>
    </button>

    <button
      class="lang-btn"
      :class="{ active: i18n.lang === 'en' }"
      @click="i18n.setLang('en')"
      type="button"
    >
      US <span class="code">EN</span>
    </button>
  </div>
</template>

<style scoped>
/* ✅ Default style (dashboard) */
.lang-switch {
  display: inline-flex;
  gap: 6px;
  padding: 6px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.28);
  border: 1px solid rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
  align-items: center;
}

.lang-btn {
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.75);
  font-weight: 800;
  font-size: 0.8rem;
  padding: 8px 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.5px;
  display: flex;
  gap: 6px;
  align-items: center;
}

.lang-btn:hover {
  color: white;
  background: rgba(255, 255, 255, 0.08);
}

.lang-btn.active {
  color: white;
  background: rgba(59, 130, 246, 0.25);
  border: 1px solid rgba(59, 130, 246, 0.45);
}

.code {
  opacity: 0.85;
}

/* ✅ Compact style (login) */
.lang-switch.is-compact {
  padding: 4px;
  border-radius: 12px;
}

.lang-switch.is-compact .lang-btn {
  font-size: 0.7rem;
  padding: 6px 10px;
  border-radius: 10px;
}
</style>
