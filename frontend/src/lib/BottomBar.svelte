<nav>
    <ul>
        <li>A</li>
        <li>B</li>
        <li class="add">
            <a href="/#/add">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    width="24"
                    height="24"
                    ><path
                        d="M12.75 7.75a.75.75 0 00-1.5 0v3.5h-3.5a.75.75 0 000 1.5h3.5v3.5a.75.75 0 001.5 0v-3.5h3.5a.75.75 0 000-1.5h-3.5v-3.5z"
                    /><path
                        fill-rule="evenodd"
                        d="M12 1C5.925 1 1 5.925 1 12s4.925 11 11 11 11-4.925 11-11S18.075 1 12 1zM2.5 12a9.5 9.5 0 1119 0 9.5 9.5 0 01-19 0z"
                    /></svg
                >
            </a>
        </li>
        <li>D</li>
        <li>E</li>
    </ul>
</nav>

<style lang="scss">
    nav {
        width: 100vw;
        height: 40px;
        background-color: rgb(205, 205, 205);
        position: fixed;
        bottom: 0;
        left: 0;

        ul {
            display: flex;
            flex-direction: row;
            width: 100%;
            height: 100%;
            justify-content: center;
            align-items: center;

            li {
                width: 18%;
                text-align: center;
            }

            li.add {
                height: 80px;
                width: 80px;
                background-color: rgb(205, 205, 205);
                border-radius: 40px;
                box-shadow: 0 2px 3px;
                display: flex;
                align-items: center;
                justify-content: center;

                svg {
                    width: 60%;
                    height: 60%;
                }
            }
        }
    }
</style>
