package server

const indexHTML = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>Receipt Printer</title>
<style>
  body { font-family: system-ui, sans-serif; max-width: 600px; margin: auto; padding: 20px; }
  label { display: block; margin: 8px 0 2px; font-weight: bold; }
  input, select { padding: 6px; width: 100%; box-sizing: border-box; margin-bottom: 10px; }
  input[type=radio] { width: auto; margin: 0 5px 0 15px; }
  .radio-group { display: flex; align-items: center; margin-bottom: 10px; }
  button { padding: 10px 20px; margin: 5px; cursor: pointer; }
  table { width: 100%; border-collapse: collapse; margin: 15px 0; }
  th, td { padding: 6px; text-align: left; border-bottom: 1px solid #ddd; }
  th { background: #f5f5f5; }
  .totals { font-size: 1.2em; margin: 15px 0; padding: 10px; background: #f9f9f9; }
  .totals div { margin: 3px 0; }
  .totals .grand { font-weight: bold; font-size: 1.3em; }
  .item-form { display: flex; gap: 10px; align-items: end; }
  .item-form input { width: auto; flex: 1; }
  .item-form button { height: 36px; }
</style></head>
<body>
  <nav class="no-print">
    <a href="/">New Receipt</a> | 
    <a href="/history">History</a> | 
    <a href="/register">Shop Settings</a>
  </nav>
  <h1>Print a Receipt</h1>

  <form id="receiptForm" method="POST" action="/print" target="_blank">
    <label>Store name</label>
    <input name="storeName" value="{{.StoreName}}" placeholder="Your Business Name" required>

    <label>Store address</label>
    <input name="storeAddr" value="{{.StoreAddr}}" placeholder="30500 Cathedral Street, Lodwar, Township" required>
     
    <label>Store Phone</label>
    <input name="storePhone" value="{{.StorePhone}}" placeholder="25476859....">

    <label>Tax PIN /  VAT</label>
    <input name="storeTaxID" value="{{.StoreTaxID}}" placeholder="Tax Pin / VAT (optional)" required>

    <h3>Items</h3>
    <div class="item-form">
      <input id="itemName" placeholder="e.g. Labour" style="flex:2">
      <input id="itemQty" type="number" placeholder="Qty" value="1" min="1" style="flex:1">
      <input id="itemPrice" type="number" step="0.01" placeholder="Price" style="flex:1">
      <button type="button" onclick="addItem()">+ Add</button>
    </div>

    <table id="itemsTable">
      <tr><th>Item</th><th>Qty</th><th>Price</th><th>Total</th><th></th></tr>
    </table>
    <input type="hidden" name="items" id="itemsInput">

    <h3>Payment Method</h3>
    <div class="radio-group">
      <input type="radio" name="paymentMethod" value="Cash" id="payCash" checked>
      <label for="payCash" style="display:inline;font-weight:normal">Cash</label>

      <input type="radio" name="paymentMethod" value="Card" id="payCard">
      <label for="payCard" style="display:inline;font-weight:normal">Card</label>

      <input type="radio" name="paymentMethod" value="Mobile Money" id="payMomo">
      <label for="payMomo" style="display:inline;font-weight:normal">Mobile Money</label>

      <input type="radio" name="paymentMethod" value="Bank Transfer" id="payBank">
      <label for="payBank" style="display:inline;font-weight:normal">Bank Transfer</label>
    </div>

    <label>Amount Tendered</label>
    <input name="paymentAmount" type="number" step="0.01" placeholder="0.00" required>

    <label>Tax Rate (e.g. 0.16 for 16%)</label>
    <input name="taxRate" value="0.16" required>

    <div class="totals" id="totals">
      <div>Subtotal: Ksh<span id="subtotal">0.00</span></div>
      <div>Tax: Ksh<span id="tax">0.00</span></div>
      <div class="grand">TOTAL: Ksh<span id="grandTotal">0.00</span></div>
      <div>Change: Ksh<span id="change">0.00</span></div>
    </div>

    <button type="submit">🖨 Print Receipt</button>
  </form>

  <script>
    let items = [];

    function addItem() {
      const name = document.getElementById('itemName').value.trim();
      const qty = parseInt(document.getElementById('itemQty').value) || 1;
      const price = parseFloat(document.getElementById('itemPrice').value) || 0;
      if (!name || price <= 0) return alert('Enter item name and price');

      items.push({ name, qty, price });
      document.getElementById('itemName').value = '';
      document.getElementById('itemQty').value = '1';
      document.getElementById('itemPrice').value = '';
      document.getElementById('itemName').focus();
      renderTable();
      calculate();
    }

    function removeItem(index) {
      items.splice(index, 1);
      renderTable();
      calculate();
    }

    function renderTable() {
      const tbody = document.getElementById('itemsTable');
      tbody.innerHTML = '<tr><th>Item</th><th>Qty</th><th>Price</th><th>Total</th><th></th></tr>';
      items.forEach((item, i) => {
        const total = (item.qty * item.price).toFixed(2);
        tbody.innerHTML += '<tr><td>'+item.name+'</td><td>'+item.qty+'</td><td>Ksh'+item.price.toFixed(2)+'</td><td>Ksh'+total+'</td><td><button type="button" onclick="removeItem('+i+')">X</button></td></tr>';
      });
      document.getElementById('itemsInput').value = items.map(i => i.name+','+i.qty+','+i.price.toFixed(2)).join(';');
    }

    function calculate() {
      const taxRate = parseFloat(document.querySelector('[name="taxRate"]').value) || 0;
      const tendered = parseFloat(document.querySelector('[name="paymentAmount"]').value) || 0;
      const subtotal = items.reduce((sum, i) => sum + i.qty * i.price, 0);
      const tax = subtotal * taxRate;
      const total = subtotal + tax;
      const change = tendered - total;

      document.getElementById('subtotal').textContent = subtotal.toFixed(2);
      document.getElementById('tax').textContent = tax.toFixed(2);
      document.getElementById('grandTotal').textContent = total.toFixed(2);
      document.getElementById('change').textContent = change >= 0 ? change.toFixed(2) : '0.00';
    }

    // Recalculate when tax rate or amount tendered changes
    document.querySelector('[name="taxRate"]').addEventListener('input', calculate);
    document.querySelector('[name="paymentAmount"]').addEventListener('input', calculate);
  </script>
</body></html>`

const registerHTML = `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>Store Settings</title>
<style>body{font-family:sans-serif;max-width:500px;margin:auto;padding:20px;}
label{display:block;margin-top:10px;} input,select{width:100%;padding:8px;margin:5px 0;}
button{margin-top:15px;padding:10px 20px;}</style></head><body>
<h1>Store Settings</h1>
<form method="POST" action="/register/save">
    <label>Business Name</label>
    <input name="storeName" value="{{.StoreName}}" placeholder="Your Business Name" required>

    <label>Address</label>
    <input name="storeAddr" value="{{.StoreAddr}}" placeholder="30500, Lodwar, Township" required>

    <label>Phone</label>
    <input name="storePhone" value="{{.StorePhone}}" placeholder="Phone (optional)">

    <label>KRA PIN / Tax ID</label>
    <input name="storeTaxID" value="{{.StoreTaxID}}" placeholder="A123456789X">

    <label>VAT Registered?</label>
    <div>
        <input type="radio" name="vatRegistered" value="true" id="vatYes" {{if .VATRegistered}}checked{{end}}>
        <label for="vatYes" style="display:inline;font-weight:normal;">Yes, I charge VAT</label>
        <input type="radio" name="vatRegistered" value="false" id="vatNo" {{if not .VATRegistered}}checked{{end}}>
        <label for="vatNo" style="display:inline;font-weight:normal;">No, I don't charge VAT</label>
    </div>

    <label>Currency</label>
    <select name="currency">
        <option value="Ksh" {{if eq .Currency "Ksh"}}selected{{end}}>Ksh (Kenya)</option>
        <option value="UGX" {{if eq .Currency "UGX"}}selected{{end}}>UGX (Uganda)</option>
        <option value="SSP" {{if eq .Currency "SSP"}}selected{{end}}>SSP (South Sudan)</option>
        <option value="USD" {{if eq .Currency "USD"}}selected{{end}}>USD</option>
    </select>

    <label>Receipt Template</label>
    <select name="templateId">
        <option value="standard" {{if eq .TemplateId "standard"}}selected{{end}}>Standard</option>
        <option value="compact" {{if eq .TemplateId "compact"}}selected{{end}}>Compact</option>
        <option value="detailed" {{if eq .TemplateId "detailed"}}selected{{end}}>Detailed</option>
    </select>

    <label>Logo Path (optional)</label>
    <input name="logoPath" value="{{.LogoPath}}" placeholder="C:\Users\Shop\logo.png">

    <button type="submit">Save Settings</button>
</form>
<p><a href="/">← Back to New Receipt</a></p>
<script>
	let itemNameInput = document.getElementById('itemName');
	let datalist = document.createElement('datalist');
	datalist.id = 'itemList';
	itemNameInput.setAttribute('list', 'itemList');
	itemNameInput.parentNode.appendChild(datalist);
	itemNameInput.addEventListener('input', function() {
		if (this.value.length >= 2) {
			fetch('/suggest?q=' + encodeURIComponent(this.value))
            		.then(r => r.json())
            		.then(items => {
                	datalist.innerHTML = '';
                	items.forEach(item => {
                    		let opt = document.createElement('option');
                    		opt.value = item.name + ' ($' + parseFloat(item.price).toFixed(2) + ')';
                    		opt.dataset.price = item.price;
                    		datalist.appendChild(opt);
                	});
            	});
    }
});
</script>
</body></html>`