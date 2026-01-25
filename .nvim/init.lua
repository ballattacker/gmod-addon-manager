vim.lsp.enable("gopls")
require("conform").formatters_by_ft.go = { "gofmt" }

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

vim.keymap.set({ "n" }, "<leader>pbl", function()
	os.exec(string.format("make linux"))
end, { desc = "build for linux" })
vim.keymap.set({ "n" }, "<leader>pbw", function()
	os.exec(string.format("make windows"))
end, { desc = "build for windows" })
