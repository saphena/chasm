# ScoreMaster v4 - Complex arithmetic

When applying complex rules to individual bonus scores (as opposed to a group score) the following is possible, starting with some definitions:-

- **BV** = points value of current bonus
- **RV** = the *results in* value of current rule
- **N** is the number of bonuses scored within the category
- **N1** is **N** - 1 
- **SV** is the resulting score


<p>If <strong>RV</strong> is 0, <strong>SV</strong> = <strong>BV</strong> * <strong>N</strong>  simple multiplication.</p>
<p>If <strong>RV</strong> is set to "multipliers", <strong>SV</strong> = <strong>BV</strong> * <strong>RV</strong> * <strong>N</strong>  simple multiplication.</p>
<p>If <strong>RV</strong> is set to "points", <strong>SV</strong> = <strong>BV</strong> * <strong>RV</strong> ^ <strong>N1</strong> exponential score.</p>
<p>So, with <strong>BV</strong> = 5 for all bonuses claimed, <strong>RV</strong> = 2, points. Successive claims give:-</p>
<ol>
	<li><strong>SV</strong> = 5 * 2 ^ 0 = 5</li>
	<li><strong>SV</strong> = 5 * 2 ^ 1 = 10</li>
	<li><strong>SV</strong> = 5 * 2 ^ 2 = 20</li>
	<li><strong>SV</strong> = 5 * 2 ^ 3 = 40</li>
</ol>
<p>So, with <strong>BV</strong> = 5, <strong>RV</strong> = 2, multipliers</p>
<ol>
	<li><strong>SV</strong> = 5 * 2 * 1 = 10</li>
	<li><strong>SV</strong> = 5 * 2 * 2 = 20</li>
	<li><strong>SV</strong> = 5 * 2 * 3 = 30</li>
	<li><strong>SV</strong> = 5 * 2 * 4 = 40</li>
</ol>
