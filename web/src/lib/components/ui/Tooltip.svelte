<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Portal } from '@ark-ui/svelte/portal';
	import { Tooltip } from '@ark-ui/svelte/tooltip';

	type Props = {
		/** What the tooltip says. */
		label: string;
		/** Which side of the element it appears on. */
		placement?: 'top' | 'bottom' | 'left' | 'right';
		/** True to leave the element alone, for a control that only needs a
		    tooltip in some states — a sidebar item while it is folded, say. */
		disabled?: boolean;
		/** The element the tooltip belongs to. It is given the props that make
		    it the trigger, which have to be spread onto it. */
		children: Snippet<[() => Record<string, unknown>]>;
	};

	let { label, placement = 'bottom', disabled = false, children }: Props = $props();
</script>

<Tooltip.Root
	{disabled}
	openDelay={250}
	closeDelay={80}
	positioning={{ placement, gutter: 6, strategy: 'fixed' }}
>
	<Tooltip.Trigger>
		{#snippet asChild(trigger)}
			{@render children(trigger as () => Record<string, unknown>)}
		{/snippet}
	</Tooltip.Trigger>

	<!-- Portalled to the end of the document, so no sticky column, scrolling
	     table or positioned heading it is opened from can cover or clip it. -->
	<Portal>
		<Tooltip.Positioner>
			<Tooltip.Content>{label}</Tooltip.Content>
		</Tooltip.Positioner>
	</Portal>
</Tooltip.Root>
