const { createApp, ref, computed, watch } = Vue;

const App = {
    setup() {
        // --- Reactive State ---
        // *** 修改点 1: 不再使用相对路径，而是动态构建基础 URL ***
        const baseUrl = `${window.location.origin}/sub`;

        const nodesYAML = ref('');
        const availableRules = ref([
            { id: 'gfw', name: 'GFW 规则', selected: true },
            { id: 'adguard', name: '屏蔽广告', selected: true },
        ]);
        const selectedRules = ref(
            availableRules.value.filter(r => r.selected).map(r => r.id)
        );
        const newChain = ref({ relay: '', landing: '' });
        const chains = ref([]);
        const subscriptionUrl = ref('');
        const errorMsg = ref('');

        // ... (其他 computed, watchers, methods 保持不变) ...
        const availableNodeNames = computed(() => {
            if (!nodesYAML.value) return [];
            try {
                const doc = jsyaml.load(nodesYAML.value);
                if (Array.isArray(doc)) {
                    return doc.map(node => node.name).filter(Boolean);
                }
                return [];
            } catch (e) {
                return [];
            }
        });

        const clashImportUrl = computed(() => {
            if (!subscriptionUrl.value) return '#';
            return `clash://install-config?url=${encodeURIComponent(subscriptionUrl.value)}`;
        });

        watch(nodesYAML, () => {
            chains.value = [];
            newChain.value = { relay: '', landing: '' };
        });

        const addChain = () => {
            if (newChain.value.relay && newChain.value.landing) {
                chains.value.push({ ...newChain.value });
                newChain.value = { relay: '', landing: '' };
            }
        };

        const removeChain = (index) => {
            chains.value.splice(index, 1);
        };

        const base64UrlEncode = (str) => {
            return btoa(unescape(encodeURIComponent(str)))
                .replace(/\+/g, '-')
                .replace(/\//g, '_')
                .replace(/=/g, '');
        };

        const generateSubscription = () => {
            errorMsg.value = '';
            subscriptionUrl.value = '';

            if (!nodesYAML.value.trim()) {
                errorMsg.value = "节点配置不能为空。";
                return;
            }
            try {
                const doc = jsyaml.load(nodesYAML.value);
                if (!Array.isArray(doc) || doc.length === 0) {
                    throw new Error('YAML is not a valid list or is empty.');
                }
                if (!doc.every(item => item && typeof item.name === 'string')) {
                    throw new Error('All nodes in the list must have a "name" field.');
                }
            } catch (e) {
                errorMsg.value = `YAML 格式错误: ${e.message}`;
                return;
            }

            const encodedNodes = base64UrlEncode(nodesYAML.value);
            const rulesString = selectedRules.value.join(',');
            const encodedChains = chains.value.length > 0 ? base64UrlEncode(JSON.stringify(chains.value)) : '';

            const params = new URLSearchParams({
                nodes: encodedNodes,
                rules: rulesString,
                client: 'meta'
            });

            if (encodedChains) {
                params.append('chains', encodedChains);
            }

            // *** 修改点 2: 使用完整的基础 URL 来拼接最终链接 ***
            subscriptionUrl.value = `${baseUrl}?${params.toString()}`;
        };

        const copyToClipboard = () => {
            if (!subscriptionUrl.value) return;

            if (!navigator.clipboard) {
                alert('您的浏览器不支持剪贴板 API，或当前环境不安全 (非 https 或 localhost)。请手动复制。');
                return;
            }

            navigator.clipboard.writeText(subscriptionUrl.value).then(() => {
                alert('订阅链接已复制到剪贴板！');
            }).catch(err => {
                console.error('复制失败:', err);
                alert(`复制失败，请手动复制。错误信息: ${err.message}`);
            });
        };

        return {
            nodesYAML,
            availableRules,
            selectedRules,
            newChain,
            chains,
            subscriptionUrl,
            errorMsg,
            availableNodeNames,
            clashImportUrl,
            addChain,
            removeChain,
            generateSubscription,
            copyToClipboard,
        };
    }
};

createApp(App).mount('#app');