<script lang="ts">
	import { PasswordInput } from '@ark-ui/svelte/password-input';
	import { RiEyeLine, RiEyeOffLine } from 'svelte-remixicon';
	import Icon from './Icon.svelte';

	type Props = {
		label: string;
		value: string;
		disabled?: boolean;
		autocomplete?: 'current-password' | 'new-password';
	};

	let {
		label,
		value = $bindable(''),
		disabled,
		autocomplete = 'current-password'
	}: Props = $props();
</script>

<PasswordInput.Root {disabled}>
	<PasswordInput.Label>{label}</PasswordInput.Label>
	<PasswordInput.Control>
		<PasswordInput.Input bind:value {autocomplete} />
		<PasswordInput.VisibilityTrigger aria-label="Show password">
			<!-- The indicator renders `children` while the password is visible
			     and `fallback` while it is hidden. -->
			<PasswordInput.Indicator>
				{#snippet fallback()}
					<Icon icon={RiEyeLine} size="1.0625rem" />
				{/snippet}
				<Icon icon={RiEyeOffLine} size="1.0625rem" />
			</PasswordInput.Indicator>
		</PasswordInput.VisibilityTrigger>
	</PasswordInput.Control>
</PasswordInput.Root>
