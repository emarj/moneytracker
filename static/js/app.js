function todayStr() {
  var t = new Date();
  return t.toISOString().split('T')[0]
}

function parseDecimal(d) {
  return (d == "") ? new Decimal(0) : new Decimal(d)
}

const defaultTx = {
  Date: todayStr(),
  Description: "",
  Amount: "",
  Shared: false,
  Shares: [],
  GeoLocation: "",
  User: {
    ID: 1
  },
  Method: {
    ID: 1
  },
  Category: {
    ID: 0
  },
  Type: {
    ID: 0
  }
};

var vm = new Vue({
  el: '#app',
  data: {
    error: "",
    edit: false,
    transaction: {},
    users: [],
    categories: [],
    types: [],
    transactions: [],
    methods: []
  },
  methods: {
    fetchInitialState: function () {
      fetch('/api/home/')
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          return res.json()
        })
        .then(data => {
          this.transactions = data.transactions;
          this.types = data.types;
          this.users = data.users;
          this.categories = data.categories;
          this.methods = data.methods;
        })
        .catch(res => this.error = res)

    },
    fetchLatest: function (n = 5, offset = 0, orderBy = "date_modified DESC, date DESC") {
      fetch('/api/transactions/?limit=' + n + '&offset=' + offset + '&orderBy=' + orderBy)
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          return res.json()
        })
        .then(data => this.transactions.push(...data))
        .catch(res => this.error = res)
    },
    reloadLatest: function () {
      n = this.transactions.length;
      this.transactions = [];
      this.fetchLatest(n);
    },
    moreLatest: function (e) {
      e.preventDefault();
      this.fetchLatest(5, this.transactions.length);
    },
    loadTx: function (uuid) {
      document.getElementById('app').scrollIntoView();
      fetch('/api/transaction/' + uuid)
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          return res.json()
        })
        .then(data => {
          this.edit = true;
          this.transaction = data;
        })
        .catch(res => this.error = res)

    },
    deleteTx: function (e) {
      e.preventDefault();

      const uuid = this.transaction.UUID;

      fetch('/api/transaction/' + uuid, {
          method: `DELETE`,
          credentials: `same-origin`
        })
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          this.reloadLatest();
          this.edit = false;
          this.resetForm();
        })
        .catch(res => this.error = res)
    },
    updateTx: function (e) {
      e.preventDefault();
      fetch('/api/transaction/', {
          method: `PUT`,
          credentials: `same-origin`,
          body: JSON.stringify(this.transaction)
        })
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          this.reloadLatest();
          this.edit = false;
          this.resetForm();
        })
        .catch(res => this.error = res)

    },
    addTx: function (e) {
      e.preventDefault();
      fetch('/api/transaction/', {
          method: `POST`,
          credentials: `same-origin`,
          body: JSON.stringify(this.transaction)
        })
        .then(res => {
          if (!res.ok) {
            throw Error(res.statusText);
          }
          this.reloadLatest();
          this.resetForm();
        })
        .catch(res => this.error = res)

    },
    resetForm: function () {
      this.transaction = JSON.parse(JSON.stringify(defaultTx)); //This is an hack
    }

  },
  created: function () {
    this.fetchInitialState();
    this.resetForm();
  }
})

Vue.component('tx-latest', {
  props: [
    'transactions'
  ],
  template: `<div class="row justify-content-center">
  <div class="col-lg-6">
    <h4 class="title mt-5">Latest entries</h4>
    <ul class="list-group list-group-flush bg-light">
      <li v-for="t in transactions" class="list-group-item expense" v-on:click="$emit('edit',t.UUID,$event)" style="cursor:pointer;">
          <strong>{{t.Date}}</strong> {{t.User.Name}} paid <strong>{{t.Amount}}&euro;</strong> by
          {{t.Method.Name}} for {{t.Description}} ({{t.Category.Name}})<span v-if="t.Shared">, and will receive back
            <span v-for="s in t.Shares">
              {{s.Quota}} from {{s.WithName}} 
            </span>
          </span>
      </li>
    </ul>
    <a v-on:click="$emit('more',$event)" href="#more">More...</a>
  </div>
</div>`
})

Vue.component('tx-form', {
  props: [
    'edit',
    'transaction',
    'categories',
    'methods',
    'types',
    'users'
  ],
  computed: {
    SharedQuota: function () {
      if (this.transaction.Shared && Array.isArray(this.transaction.Shares) && this.transaction.Shares.length > 0) {
        return this.transaction.Shares.reduce((sum, shr) => {
          const q = parseDecimal(shr.Quota);

          return q.plus(sum)
        }, new Decimal(0));
      } else return "";
    },
    OwnQuota: function () {
      if (this.transaction.Shared && Array.isArray(this.transaction.Shares) && this.transaction.Shares.length > 0) {
        const a = parseDecimal(this.transaction.Amount);
        const sq = parseDecimal(this.SharedQuota);
        return a.minus(sq).toString();
      } else return "";
    }


  },
  methods: {
  addShare: function (e) {
    e.preventDefault();
    if (!this.transaction.Shared) {
      return;
    }

    if (!Array.isArray(this.transaction.Shares)) {
      this.transaction.Shares = new Array();
    }

    this.transaction.Shares.push({
      WithID: 0,
      Quota: 0
    });
  }},
  template: `          <form>
  <fieldset>

    <div v-if="edit">
      <!-- UUID -->
      <div class="input-group input-group-lg mb-3">
        <div class="input-group-prepend">
          <label class="input-group-text" for="UUID">UUID</label>
        </div>
        <input type="text" class="form-control" name="UUID" v-bind:value="transaction.UUID" required readonly />
      </div>

      <!-- Date Created-->
      <div class="input-group input-group-lg mb-3">
        <div class="input-group-prepend">
          <label class="input-group-text" for="DateCreated">Created on</label>
        </div>
        <input class="form-control" type="datetime-local" name="DateCreated" v-bind:value="transaction.DateCreated"
          readonly />
      </div>

      <!-- Date Modified-->
      <div class="input-group input-group-lg mb-3">
        <div class="input-group-prepend">
          <label class="input-group-text" for="DateModified">Modified on</label>
        </div>
        <input class="form-control" type="datetime-local" name="DateModified" v-bind:value="transaction.DateModified"
          readonly />
      </div>
    </div>



    <!-- Type -->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="Type">Type</label>
      </div>
      <select class="custom-select" v-model="transaction.Type.ID">
        <option v-for="t in types" v-bind:value="t.ID">{{t.Name}}</option>
      </select>
    </div>




    <!-- Date -->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="Date">When</label>
      </div>
      <input class="form-control" type="date" name="Date" v-model="transaction.Date" />
    </div>

    <!-- User -->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="User">User</label>
      </div>
      <select class="custom-select" v-model="transaction.User.ID">
        <option v-for="u in users" v-bind:value="u.ID">{{u.Name}}</option>
      </select>
    </div>


    <!-- Amount -->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="Amount">Amount</label>
      </div>

      <input class="form-control" type="number" step="0.01" name="Amount" placeholder="0.00" v-model="transaction.Amount"
        required />

      <div class="input-group-append">
        <span class="input-group-text">€</span>
      </div>


    </div>

    <!-- Description -->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="Description">Description</label>
      </div>
      <input type="text" class="form-control" name="Description" placeholder="Something..." v-model="transaction.Description"
        required />

    </div>

    <!-- Payment Method -->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="MethodID">Method</label>
      </div>
      <select class="custom-select" name="MethodID" v-model="transaction.Method.ID">
        <option v-for="m in methods" v-bind:value="m.ID">{{m.Name}}</option>
      </select>
    </div>

    <!-- Sharing -->
    <div class="input-group input-group-lg">
      <div class="input-group-prepend">
        <label class="input-group-text" for="Shared">Shared</label>
        <div class="input-group-text">
          <input type="checkbox" v-model="transaction.Shared" />
        </div>
      </div>
      <input class="form-control" type="number" step="any" v-bind:value="OwnQuota" readonly style="color:red;" />
      <input class="form-control" type="number" step="any" v-bind:value="SharedQuota" readonly style="color:green;" />
    </div>

    <div class="input-group input-group-lg mb-3" v-for="s in transaction.Shares">
      <select class="custom-select" v-model="s.WithID">
        <option v-for="u in users" v-bind:value="u.ID">{{u.Name}}</option>
      </select>
      <input class="form-control" type="number" step="any" required v-model.number="s.Quota" />
      <button class="btn btn-warning">-</button>
    </div>
    <button class="btn btn-secondary" v-on:click="addShare">Add share</button>



    <!-- Categories -->
    <div class="input-group input-group-lg mb-3 mt-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="CategoryID">Category</label>
      </div>
      <select class="custom-select" name="CategoryID" v-model="transaction.Category.ID">
        <option v-for="cat in categories" v-bind:value="cat.ID" v-bind:selected="(transaction.Category.ID == cat.ID)">{{cat.Name}}</option>
      </select>
    </div>

    <!--Geolocation-->
    <div class="input-group input-group-lg mb-3">
      <div class="input-group-prepend">
        <label class="input-group-text" for="GeoLoc">Position</label>
      </div>
      <input class="form-control" type="text" id="position" name="GeoLoc" v-bind:value="transaction.GeoLoc"
        placeholder="-" readonly />
      <div class="input-group-append">
        <button class="btn btn-outline-secondary" type="button">Get Position</button>
      </div>
    </div>



    <!-- Buttons -->
    <div class="form-group mb-3 text-right">


      <button class="btn btn-danger" v-on:click="$emit('delete',$event)" v-if="edit">Delete</button>

      <button type="reset" class="btn btn-secondary">Cancel</button>

      <button type="submit" class="btn btn-primary" v-on:click="$emit('put',$event)" v-if="edit">Update</button>
      <button type="submit" class="btn btn-primary" v-on:click="$emit('post',$event)" v-if="!edit">Add</button>



    </div>


  </fieldset>
</form>
</div>

</div>`
})