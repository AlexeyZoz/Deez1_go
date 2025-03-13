<script lang="ts">
	import "../style.css";
	import { GetOS, GetDbs, Execute } from "../../wailsjs/go/main/App.js";
	import LoadScreen from "../components/ui/LoadScreen.svelte";
	import type { Db, Table } from "../types";

	let loading = $state(false);
	let errorMsg = $state("");

	let theDb: string = $state("");
	let dbs: Db[] = $state([]);

	let theTbl: string = $state("");
	let tbls: Table[] = $state([]);

	let contentColumns: string[] = $state([]);
	let content: any[] = $state([]);

	$effect(() => {
		loading = true;

		GetOS().then((r) => console.log(r));

		GetDbs()
			.then((r: Db[]) => {
				dbs = r;
				errorMsg = "";
			})
			.catch((e: string) => {
				dbs = [];
				errorMsg = e;
			})
			.finally(() => (loading = false));
	});

	const chooseDb = (e: MouseEvent) => {
		const dbTitle = (e.target as HTMLButtonElement).textContent;

		const cmd = "SELECT name FROM sqlite_master WHERE type='table'";

		loading = true;

		Execute(dbTitle as string, cmd)
			.then((r) => {
				tbls = r as Table[];
				errorMsg = "";
				theDb = dbTitle as string;
			})
			.catch((e: string) => {
				tbls = [];
				errorMsg = e;
			})
			.finally(() => (loading = false));
	};

	const chooseTable = (e: MouseEvent) => {
		const tableTitle = (e.target as HTMLButtonElement).textContent;

		const cmd = `SELECT * FROM ${tableTitle}`;

		loading = true;

		Execute(theDb, cmd)
			.then((r) => {
				content = r;
                if( r.length > 0)
				    contentColumns = Object.keys(r[0]).reverse();
				errorMsg = "";
			})
			.catch((e) => {
				content = [];
				contentColumns = [];
				errorMsg = e;
			})
			.finally(() => (loading = false));
	};

	const query = (e: SubmitEvent) => {
		e.preventDefault();

		if (theDb == "") {
			errorMsg = "No DB selected";
			return;
		}

		const inputs = new FormData(e.target as HTMLFormElement);
		const cmd = inputs.get("queryTa")?.toString() ?? "";

		loading = true;

		Execute(theDb, cmd)
			.then((r) => {
				content = r;
                console.log(r)
                if( r.length > 0)
				    contentColumns = Object.keys(r[0]).reverse();
				errorMsg = "";
			})
			.catch((e) => {
				content = [];
				contentColumns = [];
				errorMsg = e;
			})
			.finally(() => (loading = false));
	};
</script>

<main class="flex flex-col h-screen w-full overflow-hidden">
	<div class="flex w-full bg-sky-900 justify-center">
		<!-- Databases -->
		{#if dbs.length > 0}
			{#each dbs as db}
				<button
					class="flex flex-col bg-sky-600 rounded-md m-3 px-3 py-2 hover:bg-sky-800 cursor-pointer"
					onclick={chooseDb}
				>
					<span>{db.name}</span>
				</button>
			{/each}
		{:else}
			Loading D1 Local Databases...
		{/if}
	</div>
	<div class="flex w-full h-full max-md:flex-col">
		<!-- Tables -->
		{#if tbls.length > 0}
			<div
				class="w-auto flex md:flex-col max-md:flex-wrap md:h-full gap-2 p-2 overflow-auto"
			>
				{#each tbls as t}
					<button
						class="p-4 bg-sky-700 hover:bg-sky-800 rounded-md cursor-pointer"
						onclick={chooseTable}
					>
						{t.name}
					</button>
				{/each}
			</div>
		{/if}

		<div
			class={`w-full pt-2 flex flex-col overflow-auto ${tbls.length == 0 ? "mx-2" : "me-2"}`}
		>
			<!-- Query Textare -->
			{#if theDb != ""}
				<form
					onsubmit={query}
					class="flex flex-col pb-3 w-full max-md:px-2"
				>
					<textarea
						name="queryTa"
						class="bg-sky-200 p-3 text-gray-700 rounded-t-md"
						placeholder="Insert SQLite(D1) Query here..."
					></textarea>
					<button class="bg-sky-700 rounded-b-md">Query</button>
				</form>
			{/if}

			<!-- Content -->
			{#if errorMsg != ""}
				<p class="m-4 p-8 bg-sky-600 rounded-md">{errorMsg}</p>
			{:else if content.length < 1}
				No Info
			{:else}
				<table class="rounded-md">
					<thead class="bg-sky-800 sticky top-0">
						<tr>
							{#each contentColumns as col}
								<th class="text-start px-4 py-2">{col}</th>
							{/each}
						</tr>
					</thead>
					<tbody>
                        {#each content as c, i}
                            <tr class={`bg-sky-${i % 2 == 0 ? "600" : "700"}`}>
                                {#each Object.values(c).reverse() as v}
                                    <td class="text-start px-4 py-2">{v}</td>
                                {/each}
                            </tr>
                        {/each}
					</tbody>
				</table>
				<!-- {@const fields = Array.from(new Set(content.flatMap(k => Object.keys(k)))).reverse()}
				{@const block = 100 / (fields.length + 1)}

				<div class="w-full flex mb-1 bg-sky-800 rounded-sm p-2 justify-start">
					{#each fields as f}
						<span class={`flex flex-wrap`}>{f}</span>
					{/each}
				</div>

				{#each content as c, i}
					<div class={`w-full flex flex-row my-1 bg-sky-${i % 2 == 0 ? '600' : '700'} rounded-sm p-2 justify-start`}>
						{#each Object.values(c).reverse() as v}
							<span class="flex flex-wrap">{v}</span>
						{/each}
					</div>
				{/each} -->
			{/if}
		</div>
	</div>
	<!-- <div class="input-box" id="input">
		<input autocomplete="off" bind:value={name} class="bg-white" id="name" type="text"/>
		<button class="btn" onclick={greet}>Greet</button>
	</div> -->
</main>

{#if loading}
	<LoadScreen />
{/if}
