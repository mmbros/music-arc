<?xml version="1.0" encoding="iso-8859-1"?>

<!-- ******************************************************** -->
<xsl:stylesheet
	version="1.0"
	xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:ent="http://www.mmbros.it/html/ent"
	exclude-result-prefixes="ent"
>
<!-- ******************************************************** -->
<xsl:output method="xml" omit-xml-declaration="yes" />
<!-- ******************************************************** -->
<xsl:param name="mode"></xsl:param>
<xsl:param name="id"></xsl:param>

<!-- ******************************************************** -->
<xsl:param name="test_length">yes</xsl:param>
<xsl:param name="country-list" select="document('country-list.xml')/country-list" />
<!-- ******************************************************** -->
<xsl:template match="music-arc">
	<xsl:choose>
		<xsl:when test="$mode='discography'">
			<div class="w">
				<xsl:apply-templates
					select="artist-list/artist[@id = $id]"
				        mode="discography" />
			</div>
		</xsl:when>
		<xsl:when test="$mode='album'">
			<div class="w">
				<xsl:apply-templates
					select="album-list/album[@id = $id]"
				        mode="album" />
			</div>
		</xsl:when>
		<xsl:when test="$mode='playlist'">
			<div class="w">
				<xsl:apply-templates
					select="playlist-list/playlist[@id = $id]"
				        mode="playlist" />
			</div>
		</xsl:when>
		<xsl:when test="$mode='playlist-list'">
			<xsl:apply-templates select="playlist-list" />
		</xsl:when>
		<xsl:when test="$mode='cover'">
			<xsl:apply-templates
				select="playlist-list/playlist[@id = $id]"
			        mode="cover" />
		</xsl:when>

		<xsl:when test="$mode='printAlbumCover'">
			<xsl:apply-templates
				select="album-list/album[@id = $id]"
			        mode="printAlbumCover" />
		</xsl:when>

		<xsl:when test="$mode='printAlbumsTracks'">
			<xsl:apply-templates
				select="playlist-list/playlist[@id = $id]"
			        mode="printAlbumsTracks" />
		</xsl:when>

		<xsl:when test="$mode='exportAlbum'">
			<xsl:apply-templates
				select="album-list/album[@id = $id]"
			        mode="exportAlbum" />
		</xsl:when>

		<xsl:otherwise>
			<xsl:apply-templates select="artist-list" />
		</xsl:otherwise>
	</xsl:choose>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist-list">
	<xsl:call-template name="newline" />
	<table class="w">
		<xsl:call-template name="newline" />
		<thead>
			<xsl:call-template name="newline" />
			<tr>
				<th>Artist</th>
				<th>Country</th>
				<th>#</th>
			</tr>
		</thead>
		<xsl:call-template name="newline" />
		<tbody>
			<xsl:apply-templates select="artist">
				<xsl:sort select="sortName" />
			</xsl:apply-templates>
		<xsl:call-template name="newline" />
		</tbody>
		<xsl:call-template name="newline" />
		<tfoot>
			<xsl:call-template name="newline" />
			<tr>
				<th><xsl:value-of select="count(artist)" /></th>
				<th></th>
				<th><xsl:value-of select="count(//album)" /></th>
			</tr>
		</tfoot>
	<xsl:call-template name="newline" />
	</table>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist">
	<xsl:param name="albums" select="count(/music-arc/album-list/album[artist-ref=current()/@id])" />

	<xsl:call-template name="newline" />
	<tr>
		<td>
			<xsl:choose>
				<xsl:when test="$albums = 0">
					<xsl:value-of select="sortName" />
				</xsl:when>
				<xsl:otherwise>
					<a href="javascript:d('{@id}')"><xsl:value-of select="sortName" /></a>
				</xsl:otherwise>
			</xsl:choose>
		</td>
		<td>
			<!-- <xsl:value-of select="@country" /> -->
			<xsl:value-of select="$country-list/country[@id = current()/@country]/name" />
		</td>
		<td class="r">
			<xsl:value-of select="$albums" />
		</td>
	</tr>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist" mode="discography">

	<xsl:call-template name="newline" />
	<h1><xsl:value-of select="name" /></h1>

	<xsl:call-template name="newline" />
	<ul>
		<li><b>id</b>: <xsl:value-of select="@id" /></li>
		<li><b>sortName</b>: <xsl:value-of select="sortName" /></li>
		<li><b>country</b>: <xsl:value-of select="@country" /></li>
	</ul>

	<xsl:call-template name="newline" />
	<table class="w">
		<thead>
			<xsl:call-template name="newline" />
			<tr>
				<th>img</th>
				<th>title</th>
				<th>year</th>
			</tr>
		</thead>
		<tbody>
			<xsl:apply-templates mode="discography"
				             select="/music-arc/album-list/album[artist-ref = current()/@id]">
				<xsl:sort select="date" />
				<xsl:sort select="@id" />
			</xsl:apply-templates>
		</tbody>
	</table>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="album" mode="discography">
	<xsl:call-template name="newline" />
	<tr>
		<td>
			<xsl:if test="img">
				<xsl:apply-templates select="img[1]" />
			</xsl:if>
		</td>
		<td>
			<a href="javascript:a('{@id}')"><xsl:value-of select="title" /></a>
			<xsl:for-each select="/music-arc/playlist-list/playlist[list/album-ref = current()/@id]">
				<xsl:if test="position()=1">
					<br />
				</xsl:if>
				<xsl:if test="position() &gt; 1">
					<xsl:text>, </xsl:text>
				</xsl:if>
				<xsl:value-of select="@id" />

			</xsl:for-each>

		</td>
		<td><xsl:value-of select="substring(date,1,4)" /></td>
	</tr>
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="album" mode="album">

	<h2><xsl:apply-templates select="title" /></h2>
	<xsl:call-template name="newline" />
	<xsl:if test="artist-display">
		<h3><xsl:value-of select="artist-display" /></h3>
	</xsl:if>
	<xsl:call-template name="newline" />
	<ul>
		<li><b>id</b>: <xsl:value-of select="@id" /></li>
		<li><b>artist</b>:
			<xsl:for-each select="artist-ref">
				<a href="javascript:d('{.}')"><xsl:value-of select="/music-arc/artist-list/artist[@id = current()]/name" /></a>
				<xsl:if test="position() &lt; last()">
					<xsl:text>, </xsl:text>
				</xsl:if>
			</xsl:for-each>
		</li>

		<li><b>date</b>: <xsl:value-of select="date" /></li>
		<li><b>duration</b>: <xsl:value-of select="duration" /></li>
		<li><b>action</b>: <a href="javascript:printAlbumCover('{@id}')">Print Cover</a></li>
		<li><b>action</b>: <a href="javascript:exportAlbum('{@id}')">Export</a></li>
	</ul>

	<xsl:call-template name="newline" />
	<p>
		<xsl:apply-templates select="img" />
	</p>

	<xsl:call-template name="newline" />
	<table class="w">
		<thead>
			<xsl:call-template name="newline" />
			<tr>
				<th>#</th>
				<xsl:if test="artist-ref='aavv'">
					<th>artist</th>
				</xsl:if>
				<th>title</th>
				<th>duration</th>
			</tr>
		</thead>
		<tbody>
			<xsl:apply-templates select="track-list" />
		</tbody>
	</table>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="track-list">

	<xsl:param name='colspan'>
		<xsl:choose>
			<xsl:when test="../artist-ref='aavv'">4</xsl:when>
			<xsl:otherwise>3</xsl:otherwise>
		</xsl:choose>
	</xsl:param>

	<xsl:if test="title">
		<xsl:call-template name="newline" />
		<tr>
			<td colspan="{$colspan}" class="track-list"><xsl:apply-templates select="title" /></td>
		</tr>
	</xsl:if>

	<xsl:choose>
		<xsl:when test="artist-display">
			<tr>
				<td colspan="{$colspan}" class="track-list"><xsl:apply-templates select="artist-display" /></td>
			</tr>
		</xsl:when>
		<xsl:when test="artist-ref">
			<tr>
				<td colspan="{$colspan}" class="track-list"><xsl:apply-templates select="artist-ref" /></td>
			</tr>
		</xsl:when>
	</xsl:choose>

	<xsl:apply-templates select="track" />
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="track">
	<xsl:call-template name="newline" />
	<tr>
		<td class="r"><xsl:value-of select="position()" /></td>
		<xsl:if test="../../artist-ref='aavv'">
			<td>
				<xsl:choose>
					<xsl:when test="artist-display">
						<xsl:apply-templates select="artist-display" />
					</xsl:when>
					<xsl:otherwise>
						<xsl:value-of select="artist-ref" />
					</xsl:otherwise>
				</xsl:choose>
			</td>
		</xsl:if>
		<td>
			<xsl:apply-templates select="title" />
			<xsl:apply-templates select="artist-ref[@role='featuring']" />
		</td>
		<td class="r"><xsl:value-of select="duration" /></td>
	</tr>
</xsl:template>


<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="playlist-list">
	<xsl:call-template name="newline" />
	<table class="w">
		<xsl:call-template name="newline" />
		<thead>
			<xsl:call-template name="newline" />
			<tr>
				<th>id</th>
				<th>title</th>
				<th>date</th>
			</tr>
		</thead>
		<xsl:call-template name="newline" />
		<tbody>
			<xsl:for-each select="playlist">
				<xsl:sort select="@id" />
				<tr>
					<td><a href="javascript:p('{@id}')"><xsl:value-of select="@id" /></a></td>
					<td><xsl:apply-templates select="title" /></td>
					<td><xsl:apply-templates select="date" /></td>
				</tr>
			</xsl:for-each>

		<xsl:call-template name="newline" />
		</tbody>

		<xsl:call-template name="newline" />
		<tfoot>
			<xsl:call-template name="newline" />
			<tr>
				<th><xsl:value-of select="count(playlist)" /></th>
				<th></th>
				<th></th>
			</tr>
		</tfoot>
	<xsl:call-template name="newline" />
	</table>
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="playlist" mode="playlist">

	<h2><xsl:value-of select="@id" /></h2>
	<xsl:call-template name="newline" />
	<h3><xsl:apply-templates select="title" /></h3>

	<xsl:call-template name="newline" />
	<ul>
		<li><b>date</b>: <xsl:value-of select="date" /></li>
		<li><b>creator</b>: <xsl:value-of select="creator" /></li>
		<li><a href="javascript:showCover('{@id}')">cover</a></li>
		<li><a href="javascript:printAlbumsTracks('{@id}')">Album &amp; Tracks</a></li>
	</ul>

	<xsl:apply-templates select="list" mode="playlist"/>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="list" mode="playlist">
	<xsl:call-template name="newline" />
	<table class="w">
		<thead>
			<xsl:call-template name="newline" />
			<tr>
				<th>#</th>
				<th>img</th>
				<th>artist</th>
				<th>album</th>
				<th>year</th>
			</tr>
		</thead>
		<tbody>
			<xsl:for-each select="*">
				<tr>
					<td class="r"><xsl:value-of select="position()" /></td>
					<xsl:apply-templates select="." mode="playlist"/>
				</tr>
			</xsl:for-each>
		</tbody>
	</table>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="album-ref" mode="playlist" >
	<xsl:apply-templates
		select="/music-arc/album-list/album[@id = current()]"
	        mode="playlist" />
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="album" mode="playlist" >
	<td>
		<xsl:if test="img">
			<xsl:apply-templates select="img[1]" />
		</xsl:if>
	</td>
	<td><xsl:apply-templates select="artist-ref" mode="playlist" /></td>
	<td><a href="javascript:a('{@id}')"><xsl:apply-templates select="title" /></a></td>
	<td><xsl:value-of select="substring(date,1,4)" /></td>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist-ref" mode="playlist">
	<xsl:apply-templates select="/music-arc/artist-list/artist[@id = current()]" mode="playlist" />
	<xsl:if test="position() &lt; last()">
		<xsl:text>, </xsl:text>
	</xsl:if>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist" mode="playlist" >
	<a href="javascript:d('{@id}')"><xsl:value-of select="name" /></a>
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="playlist" mode="cover">
	<xsl:call-template name="newline" />
	<div class="page">
		<xsl:call-template name="newline" />
		<table class='vcenter'><tr><td>
			<xsl:call-template name="newline" />
			<div class="cover">
				<xsl:apply-templates select="list" mode="cover" />
			<xsl:call-template name="newline" />
			</div></td></tr>
		<xsl:call-template name="newline" />
		</table>
	<xsl:call-template name="newline" />
	</div>
</xsl:template>

<!-- ******************************************************** -->

<xsl:template match="list" mode="cover">

	<xsl:call-template name="newline" />
	<table>
		<xsl:call-template name="newline" />
		<tbody>
			<xsl:for-each select="*">
				<xsl:call-template name="newline" />
				<tr>
					<td class="pos"><span><xsl:value-of select="position()" /></span></td>
					<xsl:apply-templates select="." mode="cover" />
				</tr>
			</xsl:for-each>
		<xsl:call-template name="newline" />
		</tbody>
	<xsl:call-template name="newline" />
	</table>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="album-ref" mode="cover" >
	<xsl:apply-templates
		select="/music-arc/album-list/album[@id = current()]"
	        mode="cover" />
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="album" mode="cover" >
	<td class="artist"><xsl:apply-templates select="artist-ref" mode="cover" /></td>
	<td class="title"><xsl:apply-templates select="title" /></td>
	<xsl:if test="../album[date!='']">
		<td class="date"><xsl:value-of select="substring(date,1,4)" /></td>
	</xsl:if>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist-ref" mode="cover">
	<xsl:apply-templates select="/music-arc/artist-list/artist[@id = current()]" mode="cover" />
	<xsl:if test="position() &lt; last()">
		<xsl:text>, </xsl:text>
	</xsl:if>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist" mode="cover" >
	<xsl:value-of select="name" />
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="album" mode="printAlbumCover">
	<xsl:call-template name="newline" />
	<div class="page">
		<xsl:call-template name="newline" />
		<table class='vcenter'><tr><td>
			<xsl:call-template name="newline" />
			<div class="cover">
				<h1><xsl:apply-templates select="artist-ref" mode="printAlbumCover" /></h1>
				<h2><xsl:apply-templates select="title" /></h2>
				<h3><xsl:apply-templates select="date" /></h3>

				<xsl:call-template name="newline" />
				<table>
					<xsl:call-template name="newline" />
					<tbody>
						<xsl:apply-templates select="track-list" mode="printAlbumCover" />
					<xsl:call-template name="newline" />
					</tbody>
				<xsl:call-template name="newline" />
				</table>

			<xsl:call-template name="newline" />
			</div></td></tr>
		<xsl:call-template name="newline" />
		</table>
	<xsl:call-template name="newline" />
	</div>
</xsl:template>

<!-- ******************************************************** -->
<xsl:template match="track-list" mode="printAlbumCover">

	<xsl:param name='colspan'>
		<xsl:choose>
			<xsl:when test="../artist-ref='aavv'">4</xsl:when>
			<xsl:otherwise>3</xsl:otherwise>
		</xsl:choose>
	</xsl:param>

	<xsl:if test="title">
		<xsl:call-template name="newline" />
		<tr>
			<td colspan="{$colspan}" class="track-list"><xsl:apply-templates select="title" /></td>
		</tr>
	</xsl:if>

	<xsl:choose>
		<xsl:when test="artist-display">
			<tr>
				<td colspan="{$colspan}" class="track-list"><xsl:apply-templates select="artist-display" /></td>
			</tr>
		</xsl:when>
		<xsl:when test="artist-ref">
			<tr>
				<td colspan="{$colspan}" class="track-list"><xsl:apply-templates select="artist-ref" /></td>
			</tr>
		</xsl:when>
	</xsl:choose>

	<xsl:apply-templates select="track" mode="printAlbumCover" />
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="track" mode="printAlbumCover">
	<xsl:call-template name="newline" />
	<tr>
		<td class="pos"><span><xsl:value-of select="position()" /></span></td>
		<xsl:if test="../../artist-ref='aavv'">
			<td>
				<xsl:choose>
					<xsl:when test="artist-display">
						<xsl:apply-templates select="artist-display" />
					</xsl:when>
					<xsl:otherwise>
						<xsl:value-of select="artist-ref" />
					</xsl:otherwise>
				</xsl:choose>
			</td>
		</xsl:if>
		<td class="t">
			<xsl:apply-templates select="title" />
			<xsl:apply-templates select="artist-ref[@role='featuring']" />
		</td>
		<td class="r"><xsl:value-of select="duration" /></td>
	</tr>
</xsl:template>

<!-- ******************************************************** -->
<xsl:template match="artist-ref" mode="printAlbumCover">
	<xsl:apply-templates select="/music-arc/artist-list/artist[@id = current()]" mode="printAlbumCover" />
	<xsl:if test="position() &lt; last()">
		<xsl:text>, </xsl:text>
	</xsl:if>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist" mode="printAlbumCover" >
	<xsl:value-of select="name" />
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="playlist" mode="printAlbumsTracks">
	<xsl:call-template name="newline" />
	<div class="page">
		<xsl:call-template name="newline" />
		<table class='vcenter'><tr><td>
			<xsl:call-template name="newline" />
			<div class="cover">
				<xsl:call-template name="newline" />
				<table>
					<colgroup>
						<col width="1em" />
						<col />
						<col width="0*" />
					</colgroup>
					<xsl:call-template name="newline" />
					<tbody>
						<xsl:apply-templates select="list" mode="printAlbumsTracks" />
					<xsl:call-template name="newline" />
					</tbody>
				<xsl:call-template name="newline" />
				</table>

			<xsl:call-template name="newline" />
			</div></td></tr>
		<xsl:call-template name="newline" />
		</table>
	<xsl:call-template name="newline" />
	</div>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="list" mode="printAlbumsTracks">
	<xsl:apply-templates mode="printAlbumsTracks" />
</xsl:template>

<!-- ******************************************************** -->
<xsl:template match="album-ref" mode="printAlbumsTracks" >
	<xsl:apply-templates
		select="/music-arc/album-list/album[@id = current()]"
	        mode="printAlbumsTracks">
	        <xsl:with-param name="pos" select="position()" />
	</xsl:apply-templates>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="album" mode="printAlbumsTracks" >
        <xsl:param name="pos" select="0" />

	<xsl:call-template name="newline" />

	<tr class="hdr">
<!--
		<td class="pos"><b><xsl:value-of select="$pos" /></b></td>
-->
		<td class="pos2"><xsl:value-of select="$pos" /></td>

		<xsl:call-template name="newline" />
		<td class="text">
			<span class="album-creator">
				<xsl:apply-templates select="artist-ref" />
			</span>

			<xsl:text> &#8226; </xsl:text>
			<span class="album-title">
				<xsl:apply-templates select="title" />
			</span>
			<xsl:if test="date!=''">
				<xsl:text> &#8226; </xsl:text>
				<span class="album-year">
					<xsl:value-of select="substring(date,1,4)" />
				</span>
			</xsl:if>
		</td>
		<xsl:call-template name="newline" />

		<td class="length">
			<xsl:choose>
				<xsl:when test="duration!=''">
					<xsl:value-of select="duration" />
				</xsl:when>
				<xsl:when test="$test_length='yes'">
					<span class="alert">
						<xsl:text>??:??</xsl:text>
					</span>
				</xsl:when>
			</xsl:choose>
		</td>
	</tr>
	<xsl:call-template name="newline" />
	<tr class="tlr">
		<td class="tld" colspan="3">
			<xsl:apply-templates select="track-list" mode="printAlbumsTracks" />
		</td>
	</tr>
</xsl:template>

<xsl:template match="track-list" mode="printAlbumsTracks" >
	<xsl:call-template name="newline" />
	<xsl:choose>
		<xsl:when test="title">
			<xsl:if test="position() &gt; 1"><br/></xsl:if>
			<span class="tracklist-title">
				<xsl:apply-templates select="title" />
			</span>
			<xsl:text> - </xsl:text>
		</xsl:when>
		<xsl:when test="creator">
			<xsl:if test="position() &gt; 1"><br/></xsl:if>
			<span class="tracklist-creator">
				<xsl:apply-templates select="creator" />
			</span>
			<xsl:text> - </xsl:text>
		</xsl:when>
<!--
-->
		<xsl:when test="position() &gt; 1">
			<xsl:call-template name="track-separator" />
		</xsl:when>
	</xsl:choose>

	<xsl:apply-templates select="track" mode="printAlbumsTracks" />
</xsl:template>

<!-- ******************************************************** -->

<xsl:template match="track" mode="printAlbumsTracks">
	<xsl:param name="number" select="position()+count(../preceding-sibling::tracklist/track)" />

	<xsl:if test="position() &gt; 1">
		<xsl:call-template name="track-separator" />
	</xsl:if>

	<b><xsl:value-of select="$number" /></b>

<!--
	<xsl:text disable-output-escaping="yes" >&amp;nbsp;</xsl:text>
-->
	<xsl:text disable-output-escaping="yes" >&#160;</xsl:text>

<!--
	<xsl:if test="creator!=''">
		<span class="track-creator">
			<xsl:apply-templates select="creator" />
		</span>
		<xsl:text> </xsl:text>

		<xsl:if test="feat!=''">
			<xsl:value-of select="feat/@str" />
			<xsl:text> </xsl:text>
			<span class="track-creator">
				<xsl:apply-templates select="feat" />
			</span>
			<xsl:text> </xsl:text>
		</xsl:if>

		<xsl:if test="country!=''">
			<span class="country">(<xsl:value-of select="country" />)</span>
			<xsl:text> </xsl:text>
		</xsl:if>
	</xsl:if>
-->

<!--
	<tt><xsl:apply-templates select="title" /></tt>
-->
	<span class="track-title">
		<xsl:apply-templates select="title" />
	</span>

<!--
	<xsl:if test="(feat!='') and not(creator!='')">
		<span class="track-feat">
			<xsl:text> (</xsl:text>
			<xsl:value-of select="feat/@str" />
			<xsl:text> </xsl:text>
				<xsl:apply-templates select="feat" />
			<xsl:text>)</xsl:text>
		</span>
	</xsl:if>

	<xsl:if test="$show_track_length='yes'">
		<xsl:choose>
			<xsl:when test="length!=''">
				<xsl:text disable-output-escaping="yes" >&amp;nbsp;</xsl:text>
				<span class="track-length">[<xsl:value-of select="length" />]</span>
			</xsl:when>
			<xsl:when test="$test_length='yes'">
				<xsl:text disable-output-escaping="yes" >&amp;nbsp;</xsl:text>
				<span class="track-length alert">[??:??]</span>
			</xsl:when>
		</xsl:choose>
	</xsl:if>
-->
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template match="album" mode="exportAlbum">
	<div>
		<xsl:call-template name="newline" />
		<xsl:apply-templates select="track-list" mode="exportAlbum" />
	</div>
</xsl:template>

<!-- ******************************************************** -->
<xsl:template match="track-list" mode="exportAlbum">
	<xsl:apply-templates select="track" mode="exportAlbum" />
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="track" mode="exportAlbum">
	<xsl:call-template name="newline" />

	<!-- POSITION -->
	<xsl:value-of select="position()" />
	<xsl:text>;</xsl:text>

	<!-- ARTIST -->
	<xsl:choose>
		<xsl:when test="../../artist-ref='aavv'">
			<xsl:choose>
				<xsl:when test="artist-display">
					<xsl:apply-templates select="artist-display" />
				</xsl:when>
				<xsl:otherwise>
					<xsl:apply-templates select="artist-ref" mode="exportAlbum" />
				</xsl:otherwise>
			</xsl:choose>
		</xsl:when>
		<xsl:otherwise>
			<xsl:apply-templates select="../../artist-ref" mode="exportAlbum" />
		</xsl:otherwise>
	</xsl:choose>
	<xsl:text>;</xsl:text>

	<!-- ALBUM TITLE -->
	<xsl:apply-templates select="../../title" />
	<xsl:text>;</xsl:text>

	<!-- TRACK TITLE -->
	<xsl:apply-templates select="title" />
	<xsl:text>;</xsl:text>

	<!-- DATE -->
	<!-- <xsl:apply-templates select="date" /> -->
	<xsl:value-of select="substring(../../date,1,4)" />

	<br />
</xsl:template>

<!-- ******************************************************** -->
<xsl:template match="artist-ref" mode="exportAlbum">
	<xsl:apply-templates select="/music-arc/artist-list/artist[@id = current()]" mode="exportAlbum" />
	<xsl:if test="position() &lt; last()">
		<xsl:text>, </xsl:text>
	</xsl:if>
</xsl:template>
<!-- ******************************************************** -->
<xsl:template match="artist" mode="exportAlbum" >
	<xsl:value-of select="name" />
</xsl:template>

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->
<!-- ******************************************************** -->

<xsl:template name="track-separator"> &#8226; </xsl:template>
<!--
<xsl:template name="track-separator">
	<xsl:call-template name="newline" />
	<xsl:text disable-output-escaping="yes" >&amp;bull; </xsl:text>
</xsl:template>
-->


<xsl:template match="artist-ref">
	<xsl:value-of select="/music-arc/artist-list/artist[@id = current()]/name" />
</xsl:template>

<xsl:template match="artist-ref[@role='featuring']">
	<xsl:text> </xsl:text>
	<i>
		<xsl:text>(feat. </xsl:text>
		<span>
			<xsl:value-of select="/music-arc/artist-list/artist[@id = current()]/name" />
		</span>
		<xsl:text>)</xsl:text>
	</i>
</xsl:template>

<xsl:template match="img">
	<a href="{concat('img/',.)}" target="img">
		<img src="{concat('img/',.)}" title="{.}" width="50px" height="50px" />
	</a>
</xsl:template>
<!-- ******************************************************** -->

<xsl:template name="newline">
  <xsl:text disable-output-escaping="yes"><![CDATA[
]]></xsl:text>
</xsl:template>
<!-- ******************************************************** -->

<xsl:template match="ent:*">
<!--
	<xsl:value-of disable-output-escaping="yes"
		select="concat('[&amp;',local-name(),';]')" />
-->
	<span style="color:white; background-color:red"><xsl:value-of select="concat('[ent:',local-name(),']')" /></span>
</xsl:template>

<!-- bullet = black small circle -->
<xsl:template match="ent:bull">&#8226;</xsl:template>

<!-- horizontal ellipsis = three dot leader -->
<xsl:template match="ent:hellip">&#8230;</xsl:template>

<!-- Non-breaking space -->
<xsl:template match="ent:nbsp">&amp;nbsp;</xsl:template>

<!-- ******************************************************** -->
</xsl:stylesheet>
