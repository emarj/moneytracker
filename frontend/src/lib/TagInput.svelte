<script type="ts">
    import type { Tag } from "../model";

    export let existing: Tag[] = [
        { id: 1, name: "tag1" },
        { id: 3, name: "tag3" },
    ];
    export let tags: Tag[] = [{ name: "tafas" }];

    export let allowNew = true;

    export let allTags = [];
    export let newTags = [];

    $: {
        newTags = tags.filter((t) => !t.id).map((t) => t.name);
    }
    $: {
        allTags = tags.map((t) => t.name);
    }

    let tagName = "";

    const deleteTag = (name: string) => {
        tags = tags.filter((t) => t.name !== name);
    };

    const onKeyUp = (e) => {
        let exist = false;
        const key = e.key;

        if (key == " " || key == ",") {
            tagName = tagName.slice(0, -1);

            tagName = tagName.toLowerCase();

            if (tagName === "") return;
            const foundIndex = tags.findIndex((t) => t.name == tagName);
            const found = foundIndex !== -1;
            //console.log("tag found: ", found, foundIndex);

            if (!found) {
                const tag = existing.find((t) => t.name == tagName);
                const exist = tag != null;
                //console.log("tag exists:", exist);
                if (exist) {
                    tags = [...tags, tag];
                } else if (allowNew) {
                    //console.log("adding tag");
                    tags = [...tags, { name: tagName }];
                }
            }

            tagName = "";
        }
    };
    const onKeyDown = (e) => {
        const key = e.key;
        if (key == "Backspace" && tagName == "") {
            e.preventDefault();
            tags = tags.slice(0, -1);
        }
    };
</script>

<div>
    {#each tags as tag (tag.name)}
        <span class="tag" class:exist={!!tag.id}
            >{tag.name}
            <button
                class="delete"
                on:click={() => {
                    deleteTag(tag.name);
                }}>x</button
            >
        </span>
    {/each}
    <input
        type="text"
        on:keyup={onKeyUp}
        on:keydown={onKeyDown}
        bind:value={tagName}
    />
</div>

<style lang="scss">
    div {
        height: 2rem;
        width: 100%;
        background-color: white;
        border: 1px solid black;
        display: flex;

        span.tag {
            padding: 1px 5px;
            margin: 0.2rem;
            border-radius: 0.4rem;
            background-color: bisque;

            &.exist {
                background-color: cadetblue;
            }

            button.delete {
                border-left: 1px solid black;
                height: 100%;
                cursor: pointer;
            }
        }

        input {
            border: none;
            outline: none;
            background: none;
            flex-grow: 1;
        }
    }
</style>
