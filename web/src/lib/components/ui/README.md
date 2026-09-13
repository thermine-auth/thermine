# The design system

Everything the panel is built from. The point of this folder is that a page
never writes a colour, a height or a hover state — it names one.

```svelte
<Button>Save changes</Button>
<Button variant="subtle">Cancel</Button>
<Button colorPalette="danger" icon={RiDeleteBinLine}>Delete</Button>
<Button size="sm" loading={saving}>Saving…</Button>

<IconButton icon={RiRefreshLine} label="Refresh the data" />
<LinkButton href={resolve('/admin/profile')} variant="outline">Profile</LinkButton>

<Input label="email" bind:value={email} type="email" required />
<Select label="Type" bind:value={type} options={['text', 'number', 'bool']} />
<Switch label="Verified" bind:checked={verified} />
```

## Select

`options` takes bare strings, and that is all a short list of keywords needs.
Give it objects instead when the rows should carry more:

```svelte
<Select
	label="Type"
	bind:value={type}
	options={[
		{ value: 'text', label: 'Plain text', icon: RiText },
		{ value: 'bool', label: 'Bool', icon: RiToggleLine, description: 'True or false' },
		{ value: 'json', label: 'JSON', disabled: true }
	]}
	placeholder="Pick one…"
	clearable
/>
```

It is Ark UI's Select underneath, so it comes with the keyboard — arrows,
home/end, and typing a few letters to jump to a row — and with a hidden native
select for `name` and form submission. The panel is portalled, so a drawer
that scrolls or a cell that hides its overflow cannot clip it.

`hint`, `error`, `required`, `disabled` and `readOnly` behave as they do on
`Input`: a read-only select shows its value and drops the chevron rather than
disappearing.

## Three props, and what they do

Every control takes the same three, and they combine: any size, in any
variant, in any palette.

| Prop           | Values                                        | What it decides                   |
| -------------- | --------------------------------------------- | --------------------------------- |
| `size`         | `sm` `md` `lg`                                | how tall it is — 35px, 45px, 52px |
| `variant`      | `solid` `subtle` `outline` `ghost` `plain`    | how much of the palette it uses   |
| `colorPalette` | `neutral` `danger` `success` `info` `warning` | which colours those are           |

A control sets them as data attributes and stops there:

```svelte
<button class="control" data-size={size} data-variant={variant} data-palette={colorPalette}>
```

`styles/controls.css` styles `.control` by those attributes, and it never
names a colour either — only `--palette-solid`, `--palette-subtle`,
`--palette-fg` and friends. `styles/palettes.css` fills those in per palette.

That indirection is the whole trick, and it is why:

- **changing how every button looks** is one rule in `controls.css`;
- **changing what "danger" means** is one block in `palettes.css`;
- **adding a palette** is copying a block — no component changes;
- `colorPalette="danger"` recolours solid, outline and ghost correctly,
  because each variant asks the palette rather than hard-coding red.

## Where each kind of styling lives

| Looking for                                                     | It is in                |
| --------------------------------------------------------------- | ----------------------- |
| a colour, a size, a radius, a font                              | `styles/tokens.css`     |
| what `danger` or `success` means                                | `styles/palettes.css`   |
| buttons, icon buttons, link buttons                             | `styles/controls.css`   |
| inputs, selects, switches, checkboxes, drawers, tooltips, menus | `styles/ark.css`        |
| one component's own layout                                      | its own `<style>` block |

`ark.css` styles Ark UI's parts through the `data-scope` / `data-part`
attributes they render. Styling `[data-scope='field'][data-part='input']` once
is what makes every `Input` and `Textarea` in the panel match without any of
them saying so, and `[data-scope='select']` gives the select the same filled
block so it lines up with them.

## Adding a component

1. If it is shaped like a button, give it `class="control"` and pass the three
   props through. Do not write colours.
2. If it is a form field, build it on `Field.svelte` — it already lays out the
   label, the required mark, the hint and the error. A field that brings its
   own anatomy, as `Select` does, matches that layout in `ark.css` instead.
3. If it is neither, its own `<style>` block is the right place, and it should
   use tokens rather than literal values.
4. Export it from `index.ts`.
