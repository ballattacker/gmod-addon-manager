local lspconfig = require("lspconfig")
local nvlsp = require("nvchad.configs.lspconfig")

lspconfig.gopls.setup({
	on_attach = nvlsp.on_attach,
	on_init = nvlsp.on_init,
	capabilities = nvlsp.capabilities,
})

function os.exec(cmd, opts)
	opts = opts or {}
	if opts.trim == nil then
		opts.trim = false
	end
	if opts.silent == nil then
		opts.silent = true
	end

	if opts.silent then
		cmd = cmd .. " >/dev/null 2>&1"
	end

	local f = assert(io.popen(cmd, "r"))
	local s = assert(f:read("*a"))
	f:close()

	if opts.trim then
		s = string.gsub(s, "^%s+", "")
		s = string.gsub(s, "%s+$", "")
		s = string.gsub(s, "[\n\r]+", " ")
	end

	return s
end

-- vim.keymap.set({ "n" }, "<leader>pr<enter>", function()
-- 	os.exec(string.format("odin run %s -out:%s/odin", vim.fn.getcwd(), outdir))
-- end)
vim.keymap.set({ "n" }, "<leader>pb<enter>", function()
	os.exec(string.format("./build.sh"))
end)
